package dao

import (
	"github.com/orglode/navigator/model"
)

func (d *Dao) GetMachineAll() ([]model.DemoModel, error) {
	res := make([]model.DemoModel, 0)
	db := d.MySqlMaster.Table("demo_table").Find(&res)
	if db.RecordNotFound() {
		return res, nil
	}
	if db.Error != nil {
		return res, db.Error
	}
	return res, nil
}
