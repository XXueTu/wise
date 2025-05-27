package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Task TaskConfig
}

type TaskConfig struct {
	PoolSize int `json:"PoolSize"`
}
