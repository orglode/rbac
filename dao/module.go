package dao

import "github.com/orglode/navigator/model"

func (d *Dao) GetModuleAll() ([]model.Module, error) {
	res := make([]model.Module, 0)
	db := d.MySqlSlave.Table(ModuleTable).Where("status = ?", model.StatusSuccess).Find(&res)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return res, db.Error
	}
	return res, nil
}

func (d *Dao) AddModuleInfo(info model.Module) (int64, error) {
	db := d.MySqlMaster.Table(ModuleTable).Create(&info)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return 0, db.Error
	}
	return info.Id, nil
}

func (d *Dao) ModifyModuleInfo(id int64, info model.Module) (bool, error) {
	db := d.MySqlMaster.Table(ModuleTable).Where("id = ?", id).Save(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
func (d *Dao) DelModuleInfo(id int64) (bool, error) {
	res := model.Module{}
	db := d.MySqlMaster.Table(ModuleTable).Where("id = ?", id).Delete(&res)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}

func (d *Dao) AddRoleModuleInfo(info model.RoleModule) (int64, error) {
	db := d.MySqlMaster.Table(ModuleRoleTable).Create(&info)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return 0, db.Error
	}
	return info.Id, nil
}

func (d *Dao) ModifyRoleModuleInfo(id int64, info model.RoleModule) (bool, error) {
	db := d.MySqlMaster.Table(ModuleRoleTable).Where("id = ?", id).Save(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
func (d *Dao) DelRoleModuleInfo(id int64) (bool, error) {
	res := model.RoleModule{}
	db := d.MySqlMaster.Table(ModuleRoleTable).Where("id = ?", id).Delete(&res)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
