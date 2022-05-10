package middleware

const (
	MiddlewareMetricsPromHttpId = "gin-mw-metrics"
	MiddlewareMetricsKind       = "mw-kind-metrics"
)

/*
 * ErrorHandlerConfig
 */

type PromHttpMetricsHandlerConfig struct {
}

var DefaultPromHttpMetricsHandlerConfig = PromHttpMetricsHandlerConfig{}

func (h *PromHttpMetricsHandlerConfig) GetKind() string {
	return MiddlewareMetricsKind
}

type PromHttpMetricsHandlerOption func(*PromHttpMetricsHandlerConfig)
type PromHttpMetricsHandlerConfigBuilder struct {
	opts []PromHttpMetricsHandlerOption
}

/*
func (cb *PromHttpMetricsHandlerConfigBuilder) WithEndpoint(endpoint string) *PromHttpMetricsHandlerConfigBuilder {

	handlerFactoryMap := func(c *PromHttpMetricsHandlerConfig) {
		c.Endpoint = endpoint
	}

	cb.opts = append(cb.opts, handlerFactoryMap)
	return cb
}
*/

func (cb *PromHttpMetricsHandlerConfigBuilder) Build() *PromHttpMetricsHandlerConfig {
	c := DefaultPromHttpMetricsHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
