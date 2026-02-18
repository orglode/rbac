package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"

	"rbac/conf"
	"time"

	loggerV2 "github.com/orglode/hades/logger_v2"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dao struct {
	conf  *conf.Config
	db    *Mysql
	redis *redis.Client
	Rbac  *Rabc
}

func NewDao(conf *conf.Config) *Dao {
	dbMysql := &Mysql{
		mysqlMaster: initMysqlDb(conf.Db.Master),
		mysqlSlave:  initMysqlDb(conf.Db.Slave),
	}
	db := &Dao{
		conf:  conf,
		db:    dbMysql,
		redis: initRedis(conf.Redis),
	}
	db.Rbac = NewRbac(db)
	return db
}

// Transaction 事务
func (d *Dao) Transaction(ctx context.Context, fn func(tx *gorm.DB)) error {
	return nil
}

type Mysql struct {
	mysqlMaster *gorm.DB
	mysqlSlave  *gorm.DB
}

func (d *Mysql) Master(ctx context.Context) *gorm.DB {
	return d.mysqlMaster.WithContext(ctx)
}
func (d *Mysql) Slave(ctx context.Context) *gorm.DB {
	return d.mysqlSlave.WithContext(ctx)
}

// 初始化MySQL连接池
func initMysqlDb(dbConf *conf.MysqlConfig) *gorm.DB {
	db, err := sql.Open(dbConf.Drive, dbConf.Url)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(dbConf.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime) * time.Minute)
	db.SetMaxOpenConns(dbConf.MaxOpenConn)
	master, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: initDbLog(),
	})
	if err != nil {
		return nil
	}
	return master
}
func initDbLog() logger.Interface {
	gormLogger := loggerV2.NewGormLogger(loggerV2.GetSQLLogger())
	return gormLogger
}

// 初始redis
func initRedis(conf *conf.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,     // Redis 地址
		Password: conf.PassWord, // 密码
		DB:       conf.Db,       // 数据库
	})

	ctx := context.Background()

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong) // 输出: PONG
	return rdb
}

type Paging struct {
	Page int `json:"page" schema:"page"`
	Size int `json:"size" schema:"size"`
}

func (p *Paging) Offset() int {
	return (p.Page - 1) * p.Size
}
