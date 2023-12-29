package config

// options 代表配置信息加载配置项.
type options struct {
	dynamicLoader func() (string, error)
	fileLoader    func() ([]string, error)
	loaders       []Loader
}

// Files 返回配置文件列表.
func (o *options) Files() ([]string, error) {
	if o.dynamicLoader != nil {
		f, err := o.dynamicLoader()
		if err != nil {
			return nil, err
		}
		if f != "" {
			return []string{f}, nil
		}
	}
	if o.fileLoader != nil {
		return o.fileLoader()
	}
	return nil, nil
}

// Option 定义配置信息加载选项.
type Option func(*options)

// WithFileLoader 设置文件加载器.
func WithFileLoader(fl func() ([]string, error)) Option {
	return func(opts *options) {
		opts.fileLoader = fl
	}
}

// WithLoader 设置额外的配置加载器.
//
// 文件配置加载后, 应用该 Loader.
func WithLoader(loader Loader) Option {
	return func(opts *options) {
		opts.loaders = append(opts.loaders, loader)
	}
}
