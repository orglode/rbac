package dao

import (
	"github.com/orglode/navigator/model"
)

func (d *Dao) AddRoleInfo(info model.Role) (int64, error) {
	db := d.MySqlMaster.Table(roleTable).Create(&info)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return 0, db.Error
	}
	return info.Id, nil
}

func (d *Dao) ModifyRoleInfo(id int64, info model.Role) (bool, error) {
	db := d.MySqlMaster.Table(roleTable).Where("id = ?", id).Save(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
func (d *Dao) DelRoleInfo(id int64) (bool, error) {
	res := model.Role{}
	db := d.MySqlMaster.Table(roleTable).Where("id = ?", id).Delete(&res)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
