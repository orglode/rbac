package dao

import (
	"github.com/orglode/hades/logging"
	"github.com/orglode/rbac/model"
)

func (d *Dao) GetCrmRoleTypeList(p Paging) ([]model.RoleType, int64, error) {
	res := make([]model.RoleType, 0)
	var count int64
	db := d.MySqlSlave.Table(roleTypeTable)
	if err := db.Count(&count).Error; err != nil {
		logging.Errorf("err :%v", db.Error)
		return nil, 0, err
	}
	if err := db.Limit(p.Size).Offset(p.Offset()).Order("id desc").Find(&res).Error; err != nil {
		logging.Errorf("err :%v", db.Error)
		return nil, 0, err
	}
	return res, count, nil
}

func (d *Dao) GetRoleTypeAll() ([]model.RoleType, error) {
	res := make([]model.RoleType, 0)
	db := d.MySqlSlave.Table(roleTypeTable).Where("status = ?", model.StatusSuccess).Order("id desc").Find(&res)
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return res, db.Error
	}
	return res, nil
}

func (d *Dao) AddRoleTypeInfo(info model.RoleType) (int64, error) {
	db := d.MySqlMaster.Table(roleTypeTable).Create(&info)
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return 0, db.Error
	}
	return info.Id, nil
}

func (d *Dao) ModifyRoleTypeInfo(id int64, info model.RoleType) (bool, error) {
	db := d.MySqlMaster.Table(roleTypeTable).Where("id = ?", id).Updates(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
func (d *Dao) DelRoleTypeInfo(id int64) (bool, error) {
	res := model.RoleType{}
	db := d.MySqlMaster.Table(roleTypeTable).Where("id = ?", id).Delete(&res)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		logging.Errorf("err :%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
