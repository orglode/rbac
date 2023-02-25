package conf

import "github.com/BurntSushi/toml"

type Config struct {
	Env      string `toml:"env"`
	HttpAddr string `toml:"http_addr"`
	DbMaster string `toml:"db_master"`
	DbSlave  string `toml:"db_slave"`
	RedisUrl string `toml:"redis_url"`
}

func Init() *Config {
	cfg := &Config{}
	fileUrl := "config/config.toml"
	_, err := toml.DecodeFile(fileUrl, &cfg)
	if err != nil {
		panic("config.toml is err !!")
	}
	return cfg
}
