package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/orglode/navigator/api"
	"github.com/orglode/navigator/conf"
	"go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"time"
)

type Dao struct {
	conf        *conf.Config
	MySqlMaster *gorm.DB
	Redis       redis.Conn
	logger      *zap.Logger
}

func NewDao(conf *conf.Config) *Dao {
	return &Dao{
		conf:        conf,
		MySqlMaster: initMysqlDb(conf),
		Redis:       initRedis(conf),
		logger:      api.InitLogger(),
	}
}

// 初始化MySQL连接池
func initMysqlDb(conf *conf.Config) *gorm.DB {
	master, err := gorm.Open("mysql", conf.DbMaster)
	if err != nil {
		return nil
	}
	//连接池里面允许Idel的最大连接数, 这些Idel的连接 就是并发时可以同时获取的连接,也是用完后放回池里面的互用的连接, 从而提升性能。
	master.DB().SetMaxIdleConns(90)
	//设置一个连接的最长生命周期，因为数据库本身对连接有一个超时时间的设置，如果超时时间到了数据库会单方面断掉连接，此时再用连接池内的连接进行访问就会出错, 因此这个值往往要小于数据库本身的连接超时时间
	master.DB().SetConnMaxLifetime(time.Minute * 10)
	//置最大打开的连接数，默认值为0表示不限制。控制应用于数据库建立连接的数量，避免过多连接压垮数据库。
	master.DB().SetMaxOpenConns(100)
	//defer master.DB().Close()
	return master
}

// 初始redis
func initRedis(conf *conf.Config) redis.Conn {
	conn, err := redis.Dial("tcp", conf.RedisUrl)
	if err != nil {
		return nil
	}
	return conn
}
