package middleware

import (
	"github.com/gin-gonic/gin"
)

type MwHandlerRegistryConfig struct {
	ErrCfg     *ErrorHandlerConfig           `yaml:"gin-mw-error" mapstructure:"gin-mw-error" json:"gin-mw-error"`
	MetricsCfg *PromHttpMetricsHandlerConfig `yaml:"gin-mw-metrics" mapstructure:"gin-mw-metrics" json:"gin-mw-metrics"`
	TraceCfg   *TracingHandlerConfig         `yaml:"gin-mw-tracing" mapstructure:"gin-mw-tracing" json:"gin-mw-tracing"`
}

type HandlerFactory func(interface{}) MiddlewareHandler

var handlerFactoryMap = map[string]HandlerFactory{
	ErrorHandlerId:   NewErrorHandler,
	TracingHandlerId: NewTracingHandler,
	MetricsHandlerId: NewPromHttpMetricsHandler,
}

type MwHandlerRegistry map[string]gin.HandlerFunc

var registry MwHandlerRegistry = make(map[string]gin.HandlerFunc)

func InitializeHandlerRegistry(registryConfig *MwHandlerRegistryConfig) error {

	if registryConfig.ErrCfg != nil {
		registry[ErrorHandlerId] = NewErrorHandler(registryConfig.ErrCfg).HandleFunc()
	}

	if registryConfig.TraceCfg != nil {
		registry[TracingHandlerId] = NewTracingHandler(registryConfig.TraceCfg).HandleFunc()
	}

	if registryConfig.MetricsCfg != nil {
		registry[MetricsHandlerId] = NewPromHttpMetricsHandler(registryConfig.MetricsCfg).HandleFunc()
	}

	/*
		for n, i := range registryConfig {
			if hanlderFactory, ok := handlerFactoryMap[n]; ok {
				registry[n] = hanlderFactory(i).HandleFunc()
			} else {
				err := errors.New("cannot find factory for middleware handler of id: " + n)
				log.Error().Err(err).Send()
				return err
			}
		}
	*/

	return nil
}

func GetHandlerFunc(name string) gin.HandlerFunc {
	return registry[name]
}
