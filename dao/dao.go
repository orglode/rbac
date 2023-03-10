package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/orglode/rbac/conf"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Dao struct {
	conf        *conf.Config
	MySqlMaster *gorm.DB
	MySqlSlave  *gorm.DB
	Redis       redis.Conn
}

func NewDao(conf *conf.Config) *Dao {
	return &Dao{
		conf:        conf,
		MySqlMaster: initMysqlDb(conf.DbMaster),
		MySqlSlave:  initMysqlDb(conf.DbSlave),
		Redis:       initRedis(conf),
	}
}

const (
	roleTable     = "role"      //角色表
	roleTypeTable = "role_type" //角色类型

	UserTable     = "users"     //用户表
	UserRoleTable = "user_role" //用户类型

	ModuleTable     = "module"      //菜单表
	ModuleRoleTable = "role_module" //菜单角色映射表
)

type Paging struct {
	Page int `json:"page" schema:"page"`
	Size int `json:"size" schema:"size"`
}

func (p *Paging) Offset() int {
	return (p.Page - 1) * p.Size
}

type TimeRange struct {
	StartTime int64 `json:"start_time" schema:"start_time"`
	EndTime   int64 `json:"end_time" schema:"end_time"`
}

func (t TimeRange) Range(db *gorm.DB, field string) *gorm.DB {
	if t.StartTime > 0 {
		db = db.Where(field+" >= ?", time.Unix(t.StartTime, 0))
	}
	if t.EndTime > 0 {
		db = db.Where(field+" <= ?", time.Unix(t.EndTime, 0))
	}
	return db
}
