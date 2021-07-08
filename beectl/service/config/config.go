package config

import (
	"github.com/spf13/viper"
)

const (
	ApiAddr      = "api-addr"
	P2PAddr      = "p2p-addr"
	DEBUGApiAddr = "debug-api-addr"
	Pwd          = "password"
	EndPoint     = "swap-endpoint"
	Dir          = "data-dir"
)

type Config struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	ID   string `json:"id"`
	Cfg  string `json:"cfg"`
}

func (c Config) Get(key string) interface{} {
	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigFile(c.Cfg)
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	return config.Get(key)
}

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetAddr() string {
	return c.Addr
}

func (c Config) GetID() string {
	return c.ID
}

func (c Config) GetCfg() string {
	return c.Cfg
}
