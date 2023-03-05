package dao

import (
	"github.com/orglode/hades/logging"
	"github.com/orglode/rbac/model"
)

func (d *Dao) GetUserRole(uid int64) (model.UserRole, error) {
	res := model.UserRole{}
	db := d.MySqlSlave.Table(UserRoleTable).Where("user_id = ?", uid).First(&res)
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return res, db.Error
	}
	return res, nil
}

func (d *Dao) AddUserRoleInfo(info model.UserRole) (int64, error) {
	db := d.MySqlMaster.Table(UserRoleTable).Create(&info)
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return 0, db.Error
	}
	return info.Id, nil
}

func (d *Dao) ModifyUserRoleUserIdInfo(userId int64, info model.UserRole) (bool, error) {
	db := d.MySqlMaster.Table(UserRoleTable).Where("user_id = ?", userId).Updates(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return false, db.Error
	}
	return true, nil
}

func (d *Dao) ModifyUserRoleInfo(id int64, info model.UserRole) (bool, error) {
	db := d.MySqlMaster.Table(UserRoleTable).Where("id = ?", id).Updates(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return false, db.Error
	}
	return true, nil
}

func (d *Dao) DelUserRoleInfo(id int64) (bool, error) {
	res := model.UserRole{}
	db := d.MySqlMaster.Table(UserRoleTable).Where("id = ?", id).Delete(&res)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
