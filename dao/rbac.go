package dao

import (
	"context"
	"errors"
	"rbac/model"
	"time"

	"gorm.io/gorm"
)

type Rabc struct {
	db *Mysql
}

func NewRbac(db *Dao) *Rabc {
	return &Rabc{db.db}
}

func (d *Rabc) GetAccountById(ctx context.Context, id int64) (*model.Account, error) {
	var account model.Account
	err := d.db.Slave(ctx).Table(model.AccountTable).Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (d *Rabc) GetAccountList(ctx context.Context, paging Paging) ([]*model.Account, int64, error) {
	var count int64
	var result []*model.Account
	db := d.db.Slave(ctx).Table(model.AccountTable)
	// 计算总数
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	// 查询数据
	if err := db.Offset(paging.Offset()).Limit(paging.Size).Order("create_time DESC").Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

func (d *Rabc) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	now := time.Now().Unix()
	account.CreateTime = now
	account.UpdateTime = now
	err := d.db.Master(ctx).Table(model.AccountTable).Create(account).Error
	return account, err
}

func (d *Rabc) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	account.UpdateTime = time.Now().Unix()
	err := d.db.Master(ctx).Table(model.AccountTable).Where("id = ?", account.Id).Updates(account).Error
	return account, err
}

func (d *Rabc) DeleteAccount(ctx context.Context, id int64) error {
	return d.db.Master(ctx).Table(model.AccountTable).Where("id = ?", id).Delete(&model.Account{}).Error
}

func (d *Rabc) GetAccountByAccountNumber(ctx context.Context, accountNumber string) (*model.Account, error) {
	var account model.Account
	err := d.db.Slave(ctx).Table(model.AccountTable).Where("account_number = ?", accountNumber).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// ExistAccountByAccountNumber 判断是否存在相同的账号
func (d *Rabc) ExistAccountByAccountNumber(ctx context.Context, accountNumber string, id int64) (bool, error) {
	var count int64
	err := d.db.Slave(ctx).Table(model.AccountTable).Where("account_number = ? AND id != ?", accountNumber, id).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// ExistAccountByMobile 判断是否存在相同的手机号
func (d *Rabc) ExistAccountByMobile(ctx context.Context, mobile string, id int64) (bool, error) {
	var count int64
	err := d.db.Slave(ctx).Table(model.AccountTable).Where("mobile = ? and id != ?", mobile, id).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (d *Rabc) GetPageList(ctx context.Context, page Paging) ([]*model.Page, int64, error) {
	var total int64
	data := make([]*model.Page, 0)
	db := d.db.Slave(ctx).Table(model.RbacPageTable)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(page.Offset()).Limit(page.Size).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}
	if db.Error != nil {
		return nil, 0, db.Error
	}
	return data, total, nil
}

func (d *Rabc) GetAllPage(ctx context.Context) ([]*model.Page, error) {
	data := make([]*model.Page, 0)
	db := d.db.Slave(ctx).Table(model.RbacPageTable).Find(&data)
	if err := db.Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Rabc) AddPage(ctx context.Context, page *model.Page) error {
	db := d.db.Master(ctx).Table(model.RbacPageTable).Create(page)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) GetAllPageByIds(ctx context.Context, roleIds []int64) ([]*model.Page, error) {
	data := make([]*model.Page, 0)
	db := d.db.Slave(ctx).Table(model.RbacPageTable).Where("id in(?)", roleIds).Find(&data)
	if err := db.Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Rabc) UpdatePage(ctx context.Context, page *model.Page) error {
	db := d.db.Master(ctx).Table(model.RbacPageTable).Where("id=?", page.Id).Updates(page)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) DeletePage(ctx context.Context, id int64) error {
	db := d.db.Master(ctx).Table(model.RbacPageTable).Where("id=?", id).Delete(&model.Page{})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) DeletePageByParentId(ctx context.Context, parentId int64) error {
	return d.db.Master(ctx).Table(model.RbacPageTable).Where("parent_id=?", parentId).Delete(&model.Page{}).Error
}

func (d *Rabc) GetPageById(ctx context.Context, id int64) (*model.Page, error) {
	data := &model.Page{}
	err := d.db.Slave(ctx).Table(model.RbacPageTable).Where("id=?", id).First(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (d *Rabc) GetChildPages(ctx context.Context, parentId int64) ([]*model.Page, error) {
	data := make([]*model.Page, 0)
	err := d.db.Slave(ctx).Table(model.RbacPageTable).Where("parent_id=?", parentId).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetButtonByApi 按 method+path 查找按钮权限点
func (d *Rabc) GetButtonByApi(ctx context.Context, method, path string) (*model.Page, error) {
	data := &model.Page{}
	err := d.db.Slave(ctx).Table(model.RbacPageTable).
		Where("page_type=? AND api_method=? AND api_path=?", model.PageTypeButton, method, path).
		First(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// ExistRolePage 角色是否拥有某页面/按钮
func (d *Rabc) ExistRolePage(ctx context.Context, roleId, pageId int64) (bool, error) {
	var count int64
	err := d.db.Slave(ctx).Table(model.RbacRolePageTable).
		Where("role_id=? AND page_id=?", roleId, pageId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Rabc) DeleteRolePageByPageIds(ctx context.Context, pageIds []int64) error {
	if len(pageIds) == 0 {
		return nil
	}
	return d.db.Master(ctx).Table(model.RbacRolePageTable).Where("page_id in(?)", pageIds).Delete(&model.RolePage{}).Error
}

func (d *Rabc) GetAllPageRoleId(ctx context.Context, id int64) ([]*model.RolePage, error) {
	data := make([]*model.RolePage, 0)
	db := d.db.Slave(ctx).Table(model.RbacRolePageTable).Where("role_id=?", id).Find(&data)
	if err := db.Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Rabc) AddPageRole(ctx context.Context, rolePage *model.RolePage) error {
	db := d.db.Master(ctx).Table(model.RbacRolePageTable).Create(rolePage)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) DeletePageRole(ctx context.Context, roleId int64) error {
	db := d.db.Master(ctx).Table(model.RbacRolePageTable).Where("role_id=?", roleId).Delete(&model.RolePage{})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) GetRoleList(ctx context.Context, page Paging) ([]*model.Role, int64, error) {
	var total int64
	data := make([]*model.Role, 0)
	db := d.db.Slave(ctx).Table(model.RbacRoleTable)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(page.Offset()).Limit(page.Size).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (d *Rabc) GetRoleById(ctx context.Context, id int64) (*model.Role, error) {
	data := &model.Role{}
	db := d.db.Slave(ctx).Table(model.RbacRoleTable).Where("id = ?", id).First(&data)
	if err := db.Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Rabc) GetAllRole(ctx context.Context) ([]*model.Role, error) {
	data := make([]*model.Role, 0)
	db := d.db.Slave(ctx).Table(model.RbacRoleTable).Find(&data)
	if err := db.Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Rabc) AddRole(ctx context.Context, role *model.Role) error {
	db := d.db.Master(ctx).Table(model.RbacRoleTable).Create(role)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}
func (d *Rabc) UpdateRole(ctx context.Context, role *model.Role) error {
	db := d.db.Master(ctx).Table(model.RbacRoleTable).Where("id=?", role.Id).Updates(role)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) DeleteRole(ctx context.Context, id int64) error {
	db := d.db.Master(ctx).Table(model.RbacRoleTable).Where("id=?", id).Delete(&model.Role{})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) GetAccountRoleByAccountId(ctx context.Context, id int64) (*model.AccountRole, error) {
	data := &model.AccountRole{}
	db := d.db.Slave(ctx).Table(model.RbacAccountRoleTable).Where("account_id=?", id).First(&data)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err := db.Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Rabc) AddAccountRole(ctx context.Context, role *model.AccountRole) error {
	db := d.db.Master(ctx).Table(model.RbacAccountRoleTable).Create(role)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) DelAccountRole(ctx context.Context, accountId int64) error {
	db := d.db.Master(ctx).Table(model.RbacAccountRoleTable).Where("account_id=? ", accountId).Delete(&model.AccountRole{})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (d *Rabc) UpdateAccountLastLoginTime(ctx context.Context, id int64, ts int64) error {
	return d.db.Master(ctx).Table(model.AccountTable).Where("id=?", id).
		Updates(map[string]interface{}{"last_login_time": ts, "update_time": ts}).Error
}

func (d *Rabc) CreateLoginLog(ctx context.Context, log *model.LoginLog) error {
	return d.db.Master(ctx).Table(model.LoginLogTable).Create(log).Error
}

func (d *Rabc) GetLoginLogList(ctx context.Context, query model.LoginLogQuery) ([]*model.LoginLog, int64, error) {
	var count int64
	result := make([]*model.LoginLog, 0)
	db := d.db.Slave(ctx).Table(model.LoginLogTable)
	if query.AccountId > 0 {
		db = db.Where("account_id=?", query.AccountId)
	}
	if query.AccountNumber != "" {
		db = db.Where("account_number=?", query.AccountNumber)
	}
	if query.Status > 0 {
		db = db.Where("status=?", query.Status)
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	page := Paging{Page: query.Page, Size: query.Size}
	if page.Page <= 0 {
		page.Page = 1
	}
	if page.Size <= 0 {
		page.Size = 20
	}
	if err := db.Offset(page.Offset()).Limit(page.Size).Order("id DESC").Find(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, count, nil
}
