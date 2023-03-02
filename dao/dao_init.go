package dao

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	mylog "github.com/orglode/navigator/api/logger"
	"github.com/orglode/navigator/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// 初始化MySQL连接池
func initMysqlDb(dbConf string) *gorm.DB {
	db, _ := sql.Open("mysql", dbConf)
	db.SetMaxIdleConns(90)
	//设置一个连接的最长生命周期，因为数据库本身对连接有一个超时时间的设置，如果超时时间到了数据库会单方面断掉连接，此时再用连接池内的连接进行访问就会出错, 因此这个值往往要小于数据库本身的连接超时时间
	db.SetConnMaxLifetime(time.Minute * 10)
	//置最大打开的连接数，默认值为0表示不限制。控制应用于数据库建立连接的数量，避免过多连接压垮数据库。
	db.SetMaxOpenConns(100)
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
	newLogger := logger.New(
		log.New(mylog.Logger, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	return newLogger
}

// 初始redis
func initRedis(conf *conf.Config) redis.Conn {
	conn, err := redis.Dial("tcp", conf.RedisUrl)
	if err != nil {
		return nil
	}
	return conn
}
