package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	apiErr "rbac/api/error"
	"rbac/api/jwt"
	"rbac/dao"
	"rbac/model"
	"sort"
	"strconv"
	"strings"
	"time"

	logger "github.com/orglode/hades/logger_v2"
	"go.uber.org/zap"
)

type RbacService struct {
	*Service
}

func NewRbacService(service *Service) *RbacService {
	return &RbacService{service}
}

func (s *RbacService) GetAccountById(ctx context.Context, id int64) (*model.Account, error) {
	data, err := s.dao.Rbac.GetAccountById(ctx, id)
	if err != nil {
		return nil, err
	}
	data.Password = ""
	return data, nil
}

func (s *RbacService) GetAccountList(ctx context.Context, req model.AccountQuery) (interface{}, error) {
	req.Normalize()
	data, total, err := s.dao.Rbac.GetAccountList(ctx, req)
	if err != nil {
		return nil, err
	}
	for _, item := range data {
		item.Password = ""
	}
	return &model.AccountListRest{
		List:  data,
		Total: total,
	}, nil
}

func (s *RbacService) AddAccount(ctx context.Context, user *model.Account) (int64, error) {
	// 查询是否存在相同的账号和手机号
	if user.AccountNumber != "" {
		exist, err := s.dao.Rbac.ExistAccountByAccountNumber(ctx, user.AccountNumber, 0)
		if err != nil {
			return 0, err
		}
		if exist {
			return 0, apiErr.RbacAccountAlreadyExists
		}
	}
	if user.Mobile != "" {
		exist, err := s.dao.Rbac.ExistAccountByMobile(ctx, user.Mobile, 0)
		if err != nil {
			return 0, err
		}
		if exist {
			return 0, apiErr.RbacMobileExists
		}
	}

	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	if user.Status == 0 {
		user.Status = model.AccountStatusEnable
	}
	result, err := s.dao.Rbac.CreateAccount(ctx, user)
	if err != nil {
		return 0, err
	}
	return result.Id, nil
}

func (s *RbacService) UpdateAccount(ctx context.Context, user *model.Account) (bool, error) {
	if user.AccountNumber != "" {
		exist, err := s.dao.Rbac.ExistAccountByAccountNumber(ctx, user.AccountNumber, user.Id)
		if err != nil {
			return false, err
		}
		if exist {
			return false, apiErr.RbacAccountAlreadyExists
		}
	}
	if user.Mobile != "" {
		exist, err := s.dao.Rbac.ExistAccountByMobile(ctx, user.Mobile, user.Id)
		if err != nil {
			return false, err
		}
		if exist {
			return false, apiErr.RbacMobileExists
		}
	}
	if user.Password != "" {
		hash := md5.Sum([]byte(user.Password))
		user.Password = hex.EncodeToString(hash[:])
	}
	_, err := s.dao.Rbac.UpdateAccount(ctx, user)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *RbacService) DeleteAccount(ctx context.Context, id int64) error {
	return s.dao.Rbac.DeleteAccount(ctx, id)
}

func (s *RbacService) LoginAccount(ctx context.Context, accountNumber, password string, meta model.LoginMeta) (*model.AccountLoginResult, error) {
	writeLog := func(accountId int64, status int, message string) {
		_ = s.dao.Rbac.CreateLoginLog(ctx, &model.LoginLog{
			AccountId:     accountId,
			AccountNumber: accountNumber,
			LoginIp:       meta.LoginIp,
			UserAgent:     meta.UserAgent,
			Status:        status,
			Message:       message,
			CreateTime:    time.Now().Unix(),
		})
	}

	user, err := s.dao.Rbac.GetAccountByAccountNumber(ctx, accountNumber)
	if err != nil {
		writeLog(0, model.LoginStatusFail, "账号不存在或密码错误")
		return nil, apiErr.LoginPassWord
	}
	hash := md5.Sum([]byte(password))
	if user.Password != hex.EncodeToString(hash[:]) {
		writeLog(user.Id, model.LoginStatusFail, "账号密码错误")
		return nil, apiErr.LoginPassWord
	}
	if user.Status == model.AccountStatusDisable {
		writeLog(user.Id, model.LoginStatusFail, "账号已禁用")
		return nil, apiErr.AccountDisabled
	}

	now := time.Now().Unix()
	_ = s.dao.Rbac.UpdateAccountLastLoginTime(ctx, user.Id, now)
	writeLog(user.Id, model.LoginStatusSuccess, "登录成功")

	token, err := jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, err
	}
	return &model.AccountLoginResult{
		Id:            user.Id,
		Token:         token,
		UserName:      user.UserName,
		Mobile:        user.Mobile,
		AccountNumber: user.AccountNumber,
		Status:        user.Status,
		LastLoginTime: now,
		CreateTime:    user.CreateTime,
	}, nil
}

func (s *RbacService) GetLoginLogList(ctx context.Context, query model.LoginLogQuery) (*model.AccountListRest, error) {
	data, total, err := s.dao.Rbac.GetLoginLogList(ctx, query)
	if err != nil {
		return nil, err
	}
	return &model.AccountListRest{List: data, Total: total}, nil
}

func (s *RbacService) GetOperationLogList(ctx context.Context, query model.OperationLogQuery) (*model.AccountListRest, error) {
	query.Normalize()
	data, total, err := s.dao.Rbac.GetOperationLogList(ctx, query)
	if err != nil {
		return nil, err
	}
	return &model.AccountListRest{List: data, Total: total}, nil
}

func operatorIDFromCtx(ctx context.Context) int64 {
	v := ctx.Value(model.CtxOperatorID)
	if v == nil {
		return 0
	}
	id, _ := v.(int64)
	return id
}

func (s *RbacService) writeOpLog(ctx context.Context, module, action int, targetType string, targetId int64, targetName, content string) {
	opID := operatorIDFromCtx(ctx)
	var opName, opAccount string
	if opID > 0 {
		if acc, err := s.dao.Rbac.GetAccountById(ctx, opID); err == nil && acc != nil {
			opName = acc.UserName
			opAccount = acc.AccountNumber
		}
	}
	_ = s.dao.Rbac.CreateOperationLog(ctx, &model.OperationLog{
		Module:          module,
		Action:          action,
		TargetType:      targetType,
		TargetId:        targetId,
		TargetName:      targetName,
		Content:         content,
		OperatorId:      opID,
		OperatorName:    opName,
		OperatorAccount: opAccount,
		CreateTime:      time.Now().Unix(),
	})
}

func (s *RbacService) GetUserRole(ctx context.Context, accountId int64) (*model.AccountInfo, error) {
	data, err := s.dao.Rbac.GetAccountRoleByAccountId(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	roleInfo, err := s.dao.Rbac.GetRoleById(ctx, data.RoleId)
	if err != nil {
		return nil, err
	}
	result := &model.AccountInfo{
		AccountId: accountId,
		RoleId:    roleInfo.Id,
		RoleName:  roleInfo.RoleName,
	}
	return result, nil
}

func (s *RbacService) GetUserRolePage(ctx context.Context, accountId int64) (interface{}, error) {
	if _, err := s.dao.Rbac.GetAccountById(ctx, accountId); err != nil {
		return nil, err
	}

	result := &model.AccountInfo{AccountId: accountId}

	accountRole, err := s.GetUserRole(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if accountRole == nil {
		return result, nil
	}
	result.RoleId = accountRole.RoleId
	result.RoleName = accountRole.RoleName

	pageIds, err := s.getRoleAllPage(ctx, accountRole.RoleId)
	if err != nil {
		return nil, err
	}
	pageData, err := s.dao.Rbac.GetAllPageByIds(ctx, pageIds)
	if err != nil {
		return nil, err
	}
	result.PageInfo = s.buildMenuTreeV2(pageData)
	return result, nil
}

func (s *RbacService) getRoleAllPage(ctx context.Context, roleId int64) ([]int64, error) {
	rolePage, err := s.dao.Rbac.GetAllPageRoleId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if rolePage == nil {
		return nil, nil
	}
	pageIds := make([]int64, 0, len(rolePage))
	for _, v := range rolePage {
		pageIds = append(pageIds, v.PageId)
	}
	return pageIds, nil
}

func (s *RbacService) GetRolePage(ctx context.Context, roleId int64) (interface{}, error) {
	// 获取角色关联的下级菜单
	result := &model.RolePageInfo{}
	// 获取角色关联的下级菜单
	pageIds, err := s.getRoleAllPage(ctx, roleId)
	if err != nil {
		return nil, err
	}
	pageData, err := s.dao.Rbac.GetAllPageByIds(ctx, pageIds)
	if err != nil {
		return nil, err
	}
	result.RoleId = roleId
	result.PageInfo = s.buildMenuTreeV2(pageData)
	return result, nil
}

func (s *RbacService) buildMenuTreeV2(menus []*model.Page) []*model.PageInfo {
	levelMap := make(map[int][]*model.Page)
	for i := range menus {
		menu := menus[i]
		levelMap[menu.PageLevel] = append(levelMap[menu.PageLevel], menu)
	}
	for level := range levelMap {
		sortPages(levelMap[level])
	}

	toInfo := func(menu *model.Page) *model.PageInfo {
		return &model.PageInfo{
			PageId:     menu.Id,
			PageLevel:  menu.PageLevel,
			PageName:   menu.PageName,
			Icon:       menu.Icon,
			PagePath:   menu.Url,
			PageType:   menu.PageType,
			ButtonCode: menu.ButtonCode,
			ApiMethod:  menu.ApiMethod,
			ApiPath:    menu.ApiPath,
			Sort:       menu.Sort,
			CreateTime: menu.CreateTime,
			ChildPage:  make([]*model.PageInfo, 0),
			Buttons:    make([]*model.PageInfo, 0),
		}
	}

	var roots []*model.PageInfo
	secondMap := make(map[int64]*model.PageInfo) // pageId -> 二级菜单节点

	for _, menu := range levelMap[1] {
		if menu.PageType == model.PageTypeButton {
			continue
		}
		roots = append(roots, toInfo(menu))
	}

	for _, menu := range levelMap[2] {
		if menu.PageType == model.PageTypeButton {
			continue
		}
		child := toInfo(menu)
		secondMap[menu.Id] = child
		for _, root := range roots {
			if menu.ParentId == root.PageId {
				root.ChildPage = append(root.ChildPage, child)
				break
			}
		}
	}

	// 三级：按钮挂到二级菜单的 Buttons
	for _, menu := range levelMap[3] {
		btn := toInfo(menu)
		if parent, ok := secondMap[menu.ParentId]; ok {
			parent.Buttons = append(parent.Buttons, btn)
		}
	}

	return roots
}

// sortPages：sort 越大越靠前；相同则 create_time 倒序
func sortPages(pages []*model.Page) {
	sort.SliceStable(pages, func(i, j int) bool {
		a, b := pages[i], pages[j]
		if a.Sort != b.Sort {
			return a.Sort > b.Sort
		}
		return a.CreateTime > b.CreateTime
	})
}

// CheckApiPermission 按按钮绑定的 API 校验权限；无对应按钮则放行（仅登录态）
func (s *RbacService) CheckApiPermission(ctx context.Context, accountId int64, method, apiPath string) error {
	account, err := s.dao.Rbac.GetAccountById(ctx, accountId)
	if err != nil {
		return err
	}
	if account.Status == model.AccountStatusDisable {
		return apiErr.AccountDisabled
	}

	button, err := s.dao.Rbac.GetButtonByApi(ctx, method, apiPath)
	if err != nil {
		return err
	}
	if button == nil {
		return nil
	}

	accountRole, err := s.GetUserRole(ctx, accountId)
	if err != nil {
		return err
	}
	if accountRole == nil {
		return apiErr.RbacInvalidPermission
	}

	ok, err := s.dao.Rbac.ExistRolePage(ctx, accountRole.RoleId, button.Id)
	if err != nil {
		return err
	}
	if !ok {
		return apiErr.RbacInvalidPermission
	}
	return nil
}

// GetMyButtons 当前账号拥有的按钮编码列表（含 page_id），供前端显隐
func (s *RbacService) GetMyButtons(ctx context.Context, accountId int64) ([]*model.PageInfo, error) {
	if _, err := s.dao.Rbac.GetAccountById(ctx, accountId); err != nil {
		return nil, err
	}

	accountRole, err := s.GetUserRole(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if accountRole == nil {
		return []*model.PageInfo{}, nil
	}
	pageIds, err := s.getRoleAllPage(ctx, accountRole.RoleId)
	if err != nil {
		return nil, err
	}
	pages, err := s.dao.Rbac.GetAllPageByIds(ctx, pageIds)
	if err != nil {
		return nil, err
	}

	buttons := make([]*model.PageInfo, 0)
	for _, p := range pages {
		if p.PageType != model.PageTypeButton {
			continue
		}
		buttons = append(buttons, &model.PageInfo{
			PageId:     p.Id,
			PageLevel:  p.PageLevel,
			PageName:   p.PageName,
			PagePath:   p.Url,
			PageType:   p.PageType,
			ButtonCode: p.ButtonCode,
			ApiMethod:  p.ApiMethod,
			ApiPath:    p.ApiPath,
		})
	}
	return buttons, nil
}

func (s *RbacService) GetAllRole(ctx context.Context) ([]*model.Role, error) {
	return s.dao.Rbac.GetAllRole(ctx)
}

func (s *RbacService) GetAllPage(ctx context.Context) ([]*model.Page, error) {
	return s.dao.Rbac.GetAllPage(ctx)
}

// ================================Rbac CURD=============================================//

func (s *RbacService) GetRoleList(ctx context.Context, page model.Paging) (*model.PageList, error) {
	result := &model.PageList{}
	data, total, err := s.dao.Rbac.GetRoleList(ctx, dao.Paging{Page: page.Page, Size: page.Size})
	if err != nil {
		return nil, err
	}
	result.Total = total
	result.List = data
	return result, nil
}

func (s *RbacService) AddRole(ctx context.Context, role *model.Role) error {
	role.CreateTime = time.Now().Unix()
	if err := s.dao.Rbac.AddRole(ctx, role); err != nil {
		return err
	}
	s.writeOpLog(ctx, model.OpModuleRole, model.OpActionCreate, "role", role.Id, role.RoleName,
		fmt.Sprintf("新增角色「%s」", role.RoleName))
	return nil
}

func (s *RbacService) ModifyRole(ctx context.Context, role *model.Role) error {
	oldName := ""
	if old, err := s.dao.Rbac.GetRoleById(ctx, role.Id); err == nil && old != nil {
		oldName = old.RoleName
	}
	if err := s.dao.Rbac.UpdateRole(ctx, role); err != nil {
		return err
	}
	content := fmt.Sprintf("修改角色「%s」", role.RoleName)
	if oldName != "" && oldName != role.RoleName {
		content = fmt.Sprintf("修改角色「%s」→「%s」", oldName, role.RoleName)
	}
	s.writeOpLog(ctx, model.OpModuleRole, model.OpActionUpdate, "role", role.Id, role.RoleName, content)
	return nil
}

func (s *RbacService) DeleteRole(ctx context.Context, roleId int64) error {
	name := ""
	if old, err := s.dao.Rbac.GetRoleById(ctx, roleId); err == nil && old != nil {
		name = old.RoleName
	}
	if err := s.dao.Rbac.DeleteRole(ctx, roleId); err != nil {
		return err
	}
	s.writeOpLog(ctx, model.OpModuleRole, model.OpActionDelete, "role", roleId, name,
		fmt.Sprintf("删除角色「%s」", name))
	return nil
}

// GetPageList 按一级菜单分页，返回 一级→二级→按钮（1:N:N）树；支持名称、API 模糊筛选
func (s *RbacService) GetPageList(ctx context.Context, query model.PageQuery) (*model.PageList, error) {
	query.Normalize()
	all, err := s.dao.Rbac.GetAllPage(ctx)
	if err != nil {
		return nil, err
	}

	filtered := filterPagesForTree(all, query.PageName, query.Api)
	tree := s.buildMenuTreeV2(filtered)

	total := int64(len(tree))
	start := (query.Page - 1) * query.Size
	if start > len(tree) {
		start = len(tree)
	}
	end := start + query.Size
	if end > len(tree) {
		end = len(tree)
	}

	return &model.PageList{
		Total: total,
		List:  tree[start:end],
	}, nil
}

// filterPagesForTree 按名称/API 模糊匹配，保留命中节点及其 1:N:N 关联上下级
func filterPagesForTree(all []*model.Page, pageName, api string) []*model.Page {
	pageName = strings.TrimSpace(pageName)
	api = strings.TrimSpace(api)
	if pageName == "" && api == "" {
		return all
	}

	byID := make(map[int64]*model.Page, len(all))
	children := make(map[int64][]*model.Page)
	for _, p := range all {
		byID[p.Id] = p
		if p.ParentId > 0 {
			children[p.ParentId] = append(children[p.ParentId], p)
		}
	}

	match := func(p *model.Page) bool {
		if pageName != "" && strings.Contains(strings.ToLower(p.PageName), strings.ToLower(pageName)) {
			return true
		}
		if api != "" {
			key := strings.ToLower(p.ApiMethod + " " + p.ApiPath)
			if strings.Contains(key, strings.ToLower(api)) ||
				strings.Contains(strings.ToLower(p.ApiPath), strings.ToLower(api)) ||
				strings.Contains(strings.ToLower(p.ApiMethod), strings.ToLower(api)) {
				return true
			}
		}
		return false
	}

	keep := make(map[int64]struct{})
	var addAncestors func(id int64)
	addAncestors = func(id int64) {
		p, ok := byID[id]
		if !ok {
			return
		}
		if _, exists := keep[id]; exists {
			return
		}
		keep[id] = struct{}{}
		if p.ParentId > 0 {
			addAncestors(p.ParentId)
		}
	}
	var addDescendants func(id int64)
	addDescendants = func(id int64) {
		for _, c := range children[id] {
			if _, exists := keep[c.Id]; exists {
				continue
			}
			keep[c.Id] = struct{}{}
			addDescendants(c.Id)
		}
	}

	for _, p := range all {
		if !match(p) {
			continue
		}
		addAncestors(p.Id)
		// 命中一级/二级：带出下级；命中按钮：已通过祖先挂上二级与一级
		if p.PageLevel == 1 || p.PageLevel == 2 {
			addDescendants(p.Id)
		}
	}

	out := make([]*model.Page, 0, len(keep))
	for _, p := range all {
		if _, ok := keep[p.Id]; ok {
			out = append(out, p)
		}
	}
	return out
}

func (s *RbacService) AddPage(ctx context.Context, page *model.Page) error {
	now := time.Now().Unix()
	page.CreateTime = now
	if page.PageType == 0 {
		page.PageType = model.PageTypeMenu
	}

	// 显式创建按钮
	if page.PageType == model.PageTypeButton {
		if page.PageLevel == 0 {
			page.PageLevel = 3
		}
		if err := s.dao.Rbac.AddPage(ctx, page); err != nil {
			return err
		}
		s.writeOpLog(ctx, model.OpModulePage, model.OpActionCreate, "page", page.Id, page.PageName,
			fmt.Sprintf("新增按钮「%s」(%s %s)", page.PageName, page.ApiMethod, page.ApiPath))
		return nil
	}

	if err := s.dao.Rbac.AddPage(ctx, page); err != nil {
		return err
	}

	// 二级菜单且传了 api_prefix：自动挂默认 CRUD 按钮
	if page.PageLevel == 2 && page.ApiPrefix != "" {
		if err := s.createDefaultCRUDButtons(ctx, page.Id, page.ApiPrefix, now); err != nil {
			return err
		}
	}
	s.writeOpLog(ctx, model.OpModulePage, model.OpActionCreate, "page", page.Id, page.PageName,
		fmt.Sprintf("新增页面「%s」路由=%s level=%d", page.PageName, page.Url, page.PageLevel))
	return nil
}

func (s *RbacService) createDefaultCRUDButtons(ctx context.Context, parentId int64, apiPrefix string, now int64) error {
	apiPrefix = strings.TrimRight(apiPrefix, "/")
	defs := []struct {
		Name   string
		Code   string
		Method string
		Path   string
	}{
		{"查询", model.ButtonQuery, "GET", apiPrefix},
		{"新增", model.ButtonCreate, "POST", apiPrefix},
		{"修改", model.ButtonUpdate, "PUT", apiPrefix + "/:id"},
		{"删除", model.ButtonDelete, "DELETE", apiPrefix + "/:id"},
	}
	for _, d := range defs {
		btn := &model.Page{
			PageName:   d.Name,
			Url:        "",
			ParentId:   parentId,
			PageLevel:  3,
			PageType:   model.PageTypeButton,
			ButtonCode: d.Code,
			ApiMethod:  d.Method,
			ApiPath:    d.Path,
			CreateTime: now,
		}
		if err := s.dao.Rbac.AddPage(ctx, btn); err != nil {
			return err
		}
	}
	return nil
}

func (s *RbacService) ModifyPage(ctx context.Context, page *model.Page) error {
	if err := s.dao.Rbac.UpdatePage(ctx, page); err != nil {
		return err
	}
	s.writeOpLog(ctx, model.OpModulePage, model.OpActionUpdate, "page", page.Id, page.PageName,
		fmt.Sprintf("修改页面「%s」路由=%s sort=%d", page.PageName, page.Url, page.Sort))
	return nil
}

func (s *RbacService) DeletePage(ctx context.Context, id int64) error {
	page, err := s.dao.Rbac.GetPageById(ctx, id)
	if err != nil {
		return err
	}
	if page == nil {
		return nil
	}

	// 收集待删 id：自身 + 子节点（按钮）+ 若删一级则含二级及其按钮
	deleteIds := []int64{id}
	children, err := s.dao.Rbac.GetChildPages(ctx, id)
	if err != nil {
		return err
	}
	for _, child := range children {
		deleteIds = append(deleteIds, child.Id)
		grand, err := s.dao.Rbac.GetChildPages(ctx, child.Id)
		if err != nil {
			return err
		}
		for _, g := range grand {
			deleteIds = append(deleteIds, g.Id)
		}
	}

	if err := s.dao.Rbac.DeleteRolePageByPageIds(ctx, deleteIds); err != nil {
		return err
	}
	for _, pid := range deleteIds {
		if err := s.dao.Rbac.DeletePage(ctx, pid); err != nil {
			return err
		}
	}
	s.writeOpLog(ctx, model.OpModulePage, model.OpActionDelete, "page", id, page.PageName,
		fmt.Sprintf("删除页面「%s」及子节点共 %d 项", page.PageName, len(deleteIds)))
	return nil
}

func (s *RbacService) AddPageRole(ctx context.Context, request *model.RolePageRequest) error {
	roleName := ""
	if role, err := s.dao.Rbac.GetRoleById(ctx, request.RoleId); err == nil && role != nil {
		roleName = role.RoleName
	}
	//删除原有的菜单配置
	err := s.dao.Rbac.DeletePageRole(ctx, request.RoleId)
	if err != nil {
		return err
	}
	pageIds := strings.Split(request.PageIds, ",")
	count := 0
	for _, pageIdStr := range pageIds {
		pageId, err := strconv.ParseInt(strings.TrimSpace(pageIdStr), 10, 64)
		if err != nil || pageId <= 0 {
			continue
		}
		data := &model.RolePage{
			RoleId: request.RoleId,
			PageId: pageId,
		}
		err = s.dao.Rbac.AddPageRole(ctx, data)
		if err != nil {
			logger.Error(ctx, "add role page error")
			continue
		}
		count++
	}
	s.writeOpLog(ctx, model.OpModuleRolePage, model.OpActionAssign, "role_page", request.RoleId, roleName,
		fmt.Sprintf("分配角色「%s」页面权限，共 %d 项，page_ids=%s", roleName, count, request.PageIds))
	return nil
}

func (s *RbacService) AddAccountRole(ctx context.Context, role *model.AccountRole) error {
	role.CreateTime = time.Now().Unix()
	accLabel, roleName, oldRoleName := "", "", ""
	if acc, err := s.dao.Rbac.GetAccountById(ctx, role.AccountId); err == nil && acc != nil {
		accLabel = acc.AccountNumber
		if acc.UserName != "" {
			accLabel = fmt.Sprintf("%s(%s)", acc.UserName, acc.AccountNumber)
		}
	}
	if r, err := s.dao.Rbac.GetRoleById(ctx, role.RoleId); err == nil && r != nil {
		roleName = r.RoleName
	}
	if old, err := s.dao.Rbac.GetAccountRoleByAccountId(ctx, role.AccountId); err == nil && old != nil && old.RoleId > 0 {
		if old.RoleId == role.RoleId {
			return nil
		}
		if r, err := s.dao.Rbac.GetRoleById(ctx, old.RoleId); err == nil && r != nil {
			oldRoleName = r.RoleName
		}
		if err := s.dao.Rbac.DelAccountRole(ctx, role.AccountId); err != nil {
			return err
		}
	}
	if err := s.dao.Rbac.AddAccountRole(ctx, role); err != nil {
		return err
	}
	content := fmt.Sprintf("账号「%s」被分配角色「%s」", accLabel, roleName)
	if oldRoleName != "" {
		content = fmt.Sprintf("账号「%s」角色由「%s」变更为「%s」", accLabel, oldRoleName, roleName)
	}
	s.writeOpLog(ctx, model.OpModuleAccount, model.OpActionAssign, "account_role", role.AccountId, accLabel, content)
	return nil
}

func (s *RbacService) AddRolePage(ctx context.Context, roleId int64, pageIds string) error {
	err := s.dao.Rbac.DeletePageRole(ctx, roleId)
	if err != nil {
		return err
	}
	if pageIds == "" {
		return nil
	}
	pageIdList := strings.Split(pageIds, ",")
	if pageIdList == nil || len(pageIdList) == 0 {
		return nil
	}

	for _, pageIdStr := range pageIdList {
		var pageId int64
		pageId, err = strconv.ParseInt(pageIdStr, 10, 64)
		if err != nil {
			fmt.Println("转换失败:", err)
			continue
		}
		temp := &model.RolePage{
			RoleId: roleId,
			PageId: pageId,
		}
		err = s.dao.Rbac.AddPageRole(ctx, temp)
		if err != nil {
			logger.Error(ctx, "AddPageRole err :", zap.Error(err))
		}
	}
	return nil
}

func (s *RbacService) DelAccountRole(ctx context.Context, accountId int64) error {
	accLabel, roleName := "", ""
	if acc, err := s.dao.Rbac.GetAccountById(ctx, accountId); err == nil && acc != nil {
		accLabel = acc.AccountNumber
		if acc.UserName != "" {
			accLabel = fmt.Sprintf("%s(%s)", acc.UserName, acc.AccountNumber)
		}
	}
	if old, err := s.dao.Rbac.GetAccountRoleByAccountId(ctx, accountId); err == nil && old != nil && old.RoleId > 0 {
		if r, err := s.dao.Rbac.GetRoleById(ctx, old.RoleId); err == nil && r != nil {
			roleName = r.RoleName
		}
	}
	if err := s.dao.Rbac.DelAccountRole(ctx, accountId); err != nil {
		return err
	}
	content := fmt.Sprintf("移除账号「%s」的角色", accLabel)
	if roleName != "" {
		content = fmt.Sprintf("移除账号「%s」的角色「%s」", accLabel, roleName)
	}
	s.writeOpLog(ctx, model.OpModuleAccount, model.OpActionRemove, "account_role", accountId, accLabel, content)
	return nil
}
