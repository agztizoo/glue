// 支持 eva 生成代码等价替换.
package di

import "go.uber.org/dig"

type evadi struct{}

// GetAppContainer 返回全局 DI 容器对象.
//
// 兼容 eva 依赖注入，用于包名等价替换.
func GetAppContainer() *evadi {
	return &evadi{}
}

func (e *evadi) MustRegister(constructor interface{}, opts ...dig.ProvideOption) {
	MustRegister(constructor, opts...)
}

func (e *evadi) Call(function interface{}, opts ...dig.InvokeOption) error {
	return Call(function, opts...)
}
