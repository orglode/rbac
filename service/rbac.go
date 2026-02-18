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

func (s *RbacService) GetAccountList(ctx context.Context, req model.Paging) (interface{}, error) {
	paging := dao.Paging{Page: req.Page, Size: req.Size}
	data, total, err := s.dao.Rbac.GetAccountList(ctx, paging)
	for _, item := range data {
		item.Password = ""
	}
	result := model.AccountListRest{
		List:  data,
		Total: total,
	}
	return &result, err
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

func (s *RbacService) LoginAccount(ctx context.Context, accountNumber, password string) (*model.AccountLoginResult, error) {
	user, err := s.dao.Rbac.GetAccountByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, apiErr.LoginPassWord
	}

	token, err := jwt.GenerateToken(user.Id)
	if err != nil {
		return nil, err
	}
	return &model.AccountLoginResult{
		Id:            user.Id,
		Token:         token,
		UserName:      user.UserName,
		Mobile:        user.Mobile,
		AccountType:   user.AccountType,
		LastLoginTime: user.LastLoginTime,
		CreateTime:    user.CreateTime,
	}, nil
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
	accountRole, err := s.GetUserRole(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if accountRole == nil {
		return nil, nil
	}
	// 获取角色关联的下级菜单
	pageIds, err := s.getRoleAllPage(ctx, accountRole.RoleId)
	if err != nil {
		return nil, err
	}
	pageData, err := s.dao.Rbac.GetAllPageByIds(ctx, pageIds)
	if err != nil {
		return nil, err
	}
	accountRole.PageInfo = s.buildMenuTreeV2(pageData)
	return accountRole, nil
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
	// 按层级分组
	levelMap := make(map[int][]*model.Page)
	for i := range menus {
		menu := menus[i]
		levelMap[menu.PageLevel] = append(levelMap[menu.PageLevel], menu)
	}

	// 从最底层开始构建
	var roots []*model.PageInfo

	// 先处理一级菜单
	for _, menu := range levelMap[1] {
		node := &model.PageInfo{
			PageId:    menu.Id,
			PageLevel: menu.PageLevel,
			PageName:  menu.PageName,
			PagePath:  menu.Url,
			ChildPage: make([]*model.PageInfo, 0),
		}
		roots = append(roots, node)
	}

	// 处理二级菜单
	for _, menu := range levelMap[2] {
		// 在roots中查找父节点
		for _, root := range roots {
			if menu.ParentId == root.PageId {
				child := &model.PageInfo{
					PageId:    menu.Id,
					PageLevel: menu.PageLevel,
					PageName:  menu.PageName,
					PagePath:  menu.Url,
					ChildPage: make([]*model.PageInfo, 0),
				}
				root.ChildPage = append(root.ChildPage, child)
				break
			}
		}
	}

	return roots
}

func (s *RbacService) GetAllRole(ctx context.Context) ([]*model.Role, error) {
	return s.dao.Rbac.GetAllRole(ctx)
}

func (s *RbacService) GetAllPage(ctx context.Context) ([]*model.Page, error) {
	return s.dao.Rbac.GetAllPage(ctx)
}

// ================================Rbac CURD=============================================//

func (s *RbacService) GetRoleList(ctx context.Context, page model.Paging) (*model.BaseCrmResponse, error) {
	result := &model.BaseCrmResponse{}
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
	return s.dao.Rbac.AddRole(ctx, role)
}

func (s *RbacService) ModifyRole(ctx context.Context, role *model.Role) error {
	if err := s.dao.Rbac.UpdateRole(ctx, role); err != nil {
		return err
	}
	return nil
}

func (s *RbacService) DeleteRole(ctx context.Context, roleId int64) error {
	return s.dao.Rbac.DeleteRole(ctx, roleId)
}

func (s *RbacService) GetPageList(ctx context.Context, page model.Paging) (*model.PageList, error) {
	result := &model.PageList{}
	data, total, err := s.dao.Rbac.GetPageList(ctx, dao.Paging{Page: page.Page, Size: page.Size})
	if err != nil {
		return nil, err
	}
	result.Total = total
	if len(data) <= 0 {
		return result, nil
	}
	pageData := make([]*model.PageListInfo, 0, len(data))
	for _, menu := range data {
		var parentName string
		if menu.ParentId > 0 {
			for _, val := range data {
				if val.Id == menu.ParentId {
					parentName = val.PageName
					break
				}
			}
		}
		temp := &model.PageListInfo{
			Page: model.Page{
				Id:         menu.Id,
				PageLevel:  menu.PageLevel,
				PageName:   menu.PageName,
				Url:        menu.Url,
				ParentId:   menu.ParentId,
				CreateTime: menu.CreateTime,
			},
			ParentName: parentName,
		}
		pageData = append(pageData, temp)
	}
	result.List = pageData
	return result, nil
}

func (s *RbacService) AddPage(ctx context.Context, page *model.Page) error {
	page.CreateTime = time.Now().Unix()
	return s.dao.Rbac.AddPage(ctx, page)
}

func (s *RbacService) ModifyPage(ctx context.Context, page *model.Page) error {
	if err := s.dao.Rbac.UpdatePage(ctx, page); err != nil {
		return err
	}
	return nil
}
func (s *RbacService) DeletePage(ctx context.Context, id int64) error {
	return s.dao.Rbac.DeletePage(ctx, id)
}

func (s *RbacService) AddPageRole(ctx context.Context, request *model.RolePageRequest) error {
	//删除原有的菜单配置
	err := s.dao.Rbac.DeletePageRole(ctx, request.RoleId)
	if err != nil {
		return err
	}
	pageIds := strings.Split(request.PageIds, ",")
	if len(pageIds) <= 0 {
		return nil
	}
	for _, pageIdStr := range pageIds {
		pageId, err := strconv.ParseInt(pageIdStr, 10, 64)
		if err != nil {
			continue
		}
		data := &model.RolePage{
			RoleId: request.RoleId,
			PageId: pageId,
		}
		err = s.dao.Rbac.AddPageRole(ctx, data)
		if err != nil {
			logger.Error(ctx, "add role page error")
		}
	}
	return nil
}

func (s *RbacService) AddAccountRole(ctx context.Context, role *model.AccountRole) error {
	role.CreateTime = time.Now().Unix()
	return s.dao.Rbac.AddAccountRole(ctx, role)
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
	return s.dao.Rbac.DelAccountRole(ctx, accountId)
}
