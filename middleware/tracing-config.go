package middleware

const (
	TracingHandlerId   = "gin-mw-tracing"
	TracingHandlerKind = "mw-kind-tracing"
)

type TracingHandlerConfig struct {
}

var DefaultTracingHandlerConfig = TracingHandlerConfig{}

func (h *TracingHandlerConfig) GetKind() string {
	return TracingHandlerKind
}

type TracingHandlerConfigOption func(*TracingHandlerConfig)
type TracingHandlerConfigBuilder struct {
	opts []TracingHandlerConfigOption
}

func (cb *TracingHandlerConfigBuilder) Build() *TracingHandlerConfig {
	c := DefaultTracingHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
