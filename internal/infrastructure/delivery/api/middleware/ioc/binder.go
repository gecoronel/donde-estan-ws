package ioc

// A Type Binder.
type Binder interface {
	// To register an instance of a type.
	ToInstance(instance interface{})
	// To register a provider of a type.
	ToProvider(provider Provider)
}

type binderImpl struct {
	context  *contextImpl
	typeName string
}

// To register an instance of a type.
// Parameters:
// 		- instance: the instance to be registered.
func (binder *binderImpl) ToInstance(instance interface{}) {
	binder.context.providers[binder.typeName] = func() interface{} {
		return instance
	}
}

// To register a provider of a type.
// Parameters:
// 		- provider: the provider to be registered.
func (binder *binderImpl) ToProvider(provider Provider) {
	binder.context.providers[binder.typeName] = provider
}
