package ioc

// NewInjector Gets Injector from context. Service Storage.
func NewInjector(context Context) Injector {
	return &injectorImpl{
		context: context,
	}
}

// Injector A service storage.
type Injector interface {
	// GetInstance Gets some instance of a type.
	GetInstance(typeName string) interface{}
	// Context Gets the context.
	Context() Context
}

type injectorImpl struct {
	context Context
}

// GetInstance Gets some instance of a type.
// Parameters:
//   - typeName: the name of the type to be obtained.
func (injector *injectorImpl) GetInstance(typeName string) interface{} {
	if provider := injector.context.GetProvider(typeName); provider != nil {
		return provider()
	}
	return nil
}

// Context Gets the context.
func (injector *injectorImpl) Context() Context {
	return injector.context
}
