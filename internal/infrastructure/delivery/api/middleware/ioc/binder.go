package ioc

// Binder A Type Binder.
type Binder interface {
	// ToInstance To register an instance of a type.
	ToInstance(instance interface{})
	// ToProvider To register a provider of a type.
	ToProvider(provider Provider)
}

type binderImpl struct {
	context  *contextImpl
	typeName string
}

// ToInstance To register an instance of a type.
// Parameters:
//   - instance: the instance to be registered.
func (binder *binderImpl) ToInstance(instance interface{}) {
	binder.context.providers[binder.typeName] = func() interface{} {
		return instance
	}
}

// ToProvider To register a provider of a type.
// Parameters:
//   - provider: the provider to be registered.
func (binder *binderImpl) ToProvider(provider Provider) {
	binder.context.providers[binder.typeName] = provider
}
