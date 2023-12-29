package di

import "go.uber.org/dig"

// registry 提供全局注册能力.
var registry = dig.New()

func MustRegister(constructor interface{}, opts ...dig.ProvideOption) {
	if err := Register(constructor, opts...); err != nil {
		panic(err)
	}
}

func RegisterWithName(constructor interface{}, name string) error {
	return Register(constructor, dig.Name(name))
}

func Register(constructor interface{}, opts ...dig.ProvideOption) error {
	return registry.Provide(constructor, opts...)
}

func Call(function interface{}, opts ...dig.InvokeOption) error {
	return registry.Invoke(function, opts...)
}

func MustCall(function interface{}, opts ...dig.InvokeOption) {
	if err := Call(function, opts...); err != nil {
		panic(err)
	}
}
