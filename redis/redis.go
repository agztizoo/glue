package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// Options 配置选项.
type Options struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	ConnTimeout  int64  `yaml:"conn_timeout"`
	ReadTimeout  int64  `yaml:"read_timeout"`
	WriteTimeout int64  `yaml:"write_timeout"`
	PoolTimeout  int64  `yaml:"pool_timeout"`
	PoolSize     int    `yaml:"pool_size"`
}

func MustNew(conf *Options) *redis.Client {
	return New(conf)
}

func New(conf *Options) *redis.Client {
	opts := &redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.DB,
		DialTimeout:  time.Duration(conf.ConnTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		PoolSize:     conf.PoolSize,
		PoolTimeout:  time.Duration(conf.PoolTimeout) * time.Second,
	}
	return redis.NewClient(opts)
}
