package dao

import "github.com/orglode/navigator/model"

func (d *Dao) GetCrmUserList(p Paging) ([]model.Users, int64, error) {
	res := make([]model.Users, 0)
	var count int64
	db := d.MySqlSlave.Table(UserTable)
	if err := db.Count(&count).Error; err != nil {
		d.logger.Sugar().Errorf("err :%v", err)
		return nil, 0, err
	}
	if err := db.Limit(p.Size).Offset(p.Offset()).Order("id desc").Find(&res).Error; err != nil {
		d.logger.Sugar().Errorf("err :%v", err)
		return nil, 0, err
	}
	return res, count, nil
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
