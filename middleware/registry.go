package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HandlerFactory func(interface{}) MiddlewareHandler

var f = map[string]HandlerFactory{
	MiddlewareErrorId:           NewErrorHandler,
	MiddlewareTracingId:         NewTracingHandler,
	MiddlewareMetricsPromHttpId: NewPromHttpMetricsHandler,
}

type MwHandlerRegistry map[string]gin.HandlerFunc
type MwHandlerRegistryConfig = map[string]interface{}

var registry MwHandlerRegistry = make(map[string]gin.HandlerFunc)

func InitializeHandlerRegistry(registryConfig map[string]interface{}) error {

	for n, i := range registryConfig {
		if f, ok := f[n]; ok {
			registry[n] = f(i).HandleFunc()
		} else {
			err := errors.New("cannot find factory for middleware handler of id: " + n)
			log.Error().Err(err).Send()
			return err
		}
	}

	return nil
}

func GetHandlerFunc(name string) gin.HandlerFunc {
	return registry[name]
}
