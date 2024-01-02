package mysql

import (
	"fmt"
	"time"

	"github.com/agztizoo/glue/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DriverName = "mysql"
	Charset    = "utf8mb4"
)

var (
	DefaultTimeout      = 100 * time.Millisecond
	DefaultReadTimeout  = 2 * time.Second
	DefaultWriteTimeout = 5 * time.Second

	// https://github.com/go-sql-driver/mysql#important-settings
	ConnMaxLifeTime = db.DefaultConnMaxLifeTime
)

// Dialector 定义字节数据库配置与方言转换函数.
func Dialector(opts *db.Options) (gorm.Dialector, error) {
	hs := make([]func(*gorm.DB) error, 0, 2)
	// 连接池配置 Hook.
	hs = append(hs, db.ConnPoolHook(int(opts.MaxIdleConns), int(opts.MaxOpenConns), ConnMaxLifeTime))

	dial := db.WithInitializeHook(dialector, hs...)
	return dial(opts)
}

func dialector(opts *db.Options) (gorm.Dialector, error) {
	dsn := toTcpDSN(opts)
	dl := mysql.New(mysql.Config{DriverName: DriverName, DSN: dsn})
	return dl, nil
}

func toTcpDSN(opts *db.Options) string {
	f := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s"
	return fmt.Sprintf(f, opts.UserName, opts.Password, opts.Host,
		opts.Port, opts.DBName, Charset,
		getTimeout(opts), getReadTimeout(opts), getWriteTimeout(opts))
}

func getTimeout(opts *db.Options) time.Duration {
	if opts.TimeoutInMills > 0 {
		return time.Duration(opts.TimeoutInMills) * time.Millisecond
	}
	return DefaultTimeout
}

func getReadTimeout(opts *db.Options) time.Duration {
	if opts.ReadTimeoutInMills > 0 {
		return time.Duration(opts.ReadTimeoutInMills) * time.Millisecond
	}
	return DefaultReadTimeout
}

func getWriteTimeout(opts *db.Options) time.Duration {
	if opts.WriteTimeoutInMills > 0 {
		return time.Duration(opts.WriteTimeoutInMills) * time.Millisecond
	}
	return DefaultWriteTimeout
}
