package dao

import (
	"github.com/orglode/navigator/model"
)

func (d *Dao) GetCrmRoleList(req model.RoleListRequest, p Paging) ([]model.RoleListInfo, int64, error) {
	res := make([]model.RoleListInfo, 0)
	var count int64
	db := d.MySqlSlave.Table(roleTable).Select("role.*,role_type.type_name").
		Joins("left join role_type on role.type_id = role_type.id")
	if req.RoleType > 0 {
		db = db.Where("role.type_id = ?", req.RoleType)
	}
	if req.Status > 0 {
		db = db.Where("role.status = ?", req.Status)
	}
	if req.RoleName != "" {
		db = db.Where("role.name = ?", req.RoleName)
	}
	if err := db.Count(&count).Error; err != nil {
		d.logger.Sugar().Errorf("err :%v", err)
		return nil, 0, err
	}
	if err := db.Limit(p.Size).Offset(p.Offset()).Order("role.id desc").Find(&res).Error; err != nil {
		d.logger.Sugar().Errorf("err :%v", err)
		return nil, 0, err
	}
	return res, count, nil
}

func (d *Dao) GetRoleAll() ([]model.Role, error) {
	res := make([]model.Role, 0)
	db := d.MySqlSlave.Table(roleTable).Where("status = ?", model.StatusSuccess).Order("id desc").Find(&res)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return res, db.Error
	}
	return res, nil
}

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
