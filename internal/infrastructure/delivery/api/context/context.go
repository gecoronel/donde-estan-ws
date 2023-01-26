package context

import (
	"context"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
)

const (
	InjectorKey         contextKey = "ioc"
	ExecutionContextKey contextKey = "execution-context"
)

type contextKey string

func SetServiceLocator(ctx context.Context, serviceLocator gateway.ServiceLocator) context.Context {
	ctx = context.WithValue(ctx, InjectorKey, serviceLocator)
	return ctx
}

func GetServiceLocator(ctx context.Context) gateway.ServiceLocator {
	return ctx.Value(InjectorKey).(gateway.ServiceLocator)
}

/*func SetExecutionContext(ctx context.Context, executionContext webapi.ExecutionContext) context.Context {
	ctx = context.WithValue(ctx, ExecutionContextKey, executionContext)
	return ctx
}

func GetExecutionContext(ctx context.Context) webapi.ExecutionContext {
	return ctx.Value(ExecutionContextKey).(webapi.ExecutionContext)
}*/

/*
const (
	ContextKey             = "{{CONTEXT}}"
	InjectorKey contextKey = "ioc"
)

type contextKey string

// GetContext Obtain the Context from the Gin Context.
func GetContext(c *gin.Context) context.Context {
	if v, ok := c.Get(ContextKey); ok {
		return v.(context.Context)
	}

	ctx := context.Background()
	c.Set(ContextKey, ctx)

	return ctx
}

// UpdateContext Update the context.
func UpdateContext(context context.Context, c *gin.Context) {
	c.Set(ContextKey, context)
}

// GetInjectorFromContext Obtain the injector (service locator) from a Context.
func GetInjectorFromContext(ctx context.Context) gateway.ServiceLocator {
	return ctx.Value(InjectorKey).(gateway.ServiceLocator)
}

// GetInjectorFromGinContext Obtain the injector (service locator) from the Gin Context.
func GetInjectorFromGinContext(c *gin.Context) gateway.ServiceLocator {
	ctx := GetContext(c)
	return GetInjectorFromContext(ctx)
}

// SetInjector Set the injector (service locator) in the Context from the Gin Context.
func SetInjector(c *gin.Context, injector gateway.ServiceLocator) {
	UpdateContext(context.WithValue(GetContext(c), InjectorKey, injector), c)
}
*/
