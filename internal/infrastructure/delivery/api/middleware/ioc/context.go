package ioc

// NewContext Creates a new Context.
func NewContext() Context {
	return &contextImpl{
		providers: map[string]Provider{},
	}
}

// Context Represents an Application Context to store some values.
type Context interface {
	// Binds some type. Registers a type.
	Bind(typeName string) Binder
	// Gets the provider to some type.
	GetProvider(typeName string) Provider
}

type contextImpl struct {
	providers map[string]Provider
}

// Bind some type.
// Parameters:
// 		- typeName: the name of the type you wish bind.
func (context *contextImpl) Bind(typeName string) Binder {
	return &binderImpl{
		context:  context,
		typeName: typeName,
	}
}

// GetProvider Gets the provider to some type.
// Parameters:
// 		- typeName: the name of the type you wish obtain a provider.
func (context *contextImpl) GetProvider(typeName string) Provider {
	if provider, ok := context.providers[typeName]; ok {
		return provider
	}
	return nil
}
