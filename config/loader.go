package config

import (
	"github.com/jinzhu/configor"
)

// Load 实现配置信息加载.
func Load(out interface{}, opts ...Option) error {
	los := &options{
		dynamicLoader: dynamicFile,
		fileLoader:    EnvAwareFile(EnvAwareFilePattern),
		loaders:       make([]Loader, 0),
	}
	for _, opt := range opts {
		opt(los)
	}
	files, err := los.Files()
	if err != nil {
		return err
	}
	if err := configor.New(&configor.Config{}).Load(out, files...); err != nil {
		return err
	}
	for _, loader := range los.loaders {
		if err := loader.Load(out); err != nil {
			return err
		}
	}
	return nil
}

// MustLoad 实现配置信息加载.
func MustLoad(out interface{}, opts ...Option) {
	if err := Load(out, opts...); err != nil {
		panic(err)
	}
}

// Loader 定义配置加载器.
type Loader interface {
	// Load 加载配置信息.
	//
	// 不同实现可以覆盖或增量加载配置.
	Load(out interface{}) error
}
