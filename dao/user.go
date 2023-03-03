package dao

import "github.com/orglode/navigator/model"

func (d *Dao) GetCrmUserList(req model.CrmUserListRequest, p Paging) ([]model.CrmUserListInfo, int64, error) {
	res := make([]model.CrmUserListInfo, 0)
	var count int64
	db := d.MySqlSlave.Table(UserTable).Select("users.*,role_type.type_name as role_type_name,role.name as role_name").
		Joins("left join user_role on user_role.user_id = users.id").
		Joins("left join role on role.id = user_role.role_id").
		Joins("left join role_type on role.type_id = role_type.id")
	if req.Phone != "" {
		db = db.Where("users.phone = ?", req.Phone)
	}
	if req.Status > 0 {
		db = db.Where("users.status = ?", req.Status)
	}
	if req.Username != "" {
		db = db.Where("users.username = ?", req.Username)
	}
	if err := db.Count(&count).Error; err != nil {
		d.logger.Sugar().Errorf("err :%v", err)
		return nil, 0, err
	}
	if err := db.Limit(p.Size).Offset(p.Offset()).Order("users.id desc").Find(&res).Error; err != nil {
		d.logger.Sugar().Errorf("err :%v", err)
		return nil, 0, err
	}
	return res, count, nil
}

func (d *Dao) GetUserInfoByAccount(account string) (model.Users, error) {
	res := model.Users{}
	db := d.MySqlSlave.Table(UserTable).Where("account = ?", account).First(&res)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err :%v", db.Error)
	}
	return res, db.Error
}

func (d *Dao) GetUserInfoById(id int64) (model.Users, error) {
	res := model.Users{}
	db := d.MySqlSlave.Table(UserTable).Where("id = ?", id).First(&res)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err :%v", db.Error)
	}
	return res, db.Error
}

func (d *Dao) AddUserInfo(info model.Users) (int64, error) {
	db := d.MySqlMaster.Table(UserTable).Create(&info)
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return 0, db.Error
	}
	return info.Id, nil
}

func (d *Dao) ModifyUserInfo(id int64, info model.Users) (bool, error) {
	db := d.MySqlMaster.Table(UserTable).Where("id = ?", id).Updates(&info)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
func (d *Dao) DelUserInfo(id int64) (bool, error) {
	res := model.Users{}
	db := d.MySqlMaster.Table(UserTable).Where("id = ?", id).Delete(&res)
	if db.RowsAffected <= 0 {
		return false, nil
	}
	if db.Error != nil {
		d.logger.Sugar().Errorf("err:%v", db.Error)
		return false, db.Error
	}
	return true, nil
}
