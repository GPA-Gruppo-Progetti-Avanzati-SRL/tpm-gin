package middleware

const (
	ErrorHandlerId              = "gin-mw-error"
	ErrorHandlerKind            = "mw-kind-error"
	ErrorHandlerDefaultWithInfo = true
)

/*
 * ErrorHandlerConfig
 */
type ErrorHandlerConfig struct {
	WithInfo bool `json:"with-info"  yaml:"with-info"  mapstructure:"with-info"`
}

var DefaultErrorHandlerConfig = ErrorHandlerConfig{
	WithInfo: ErrorHandlerDefaultWithInfo,
}

func (h *ErrorHandlerConfig) GetKind() string {
	return ErrorHandlerKind
}

type ErrorHandlerConfigOption func(*ErrorHandlerConfig)
type ErrorHandlerConfigBuilder struct {
	opts []ErrorHandlerConfigOption
}

func (cb *ErrorHandlerConfigBuilder) WithErrorEnabled(enabled bool) *ErrorHandlerConfigBuilder {

	f := func(c *ErrorHandlerConfig) {
		c.WithInfo = enabled
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *ErrorHandlerConfigBuilder) Build() *ErrorHandlerConfig {
	c := DefaultErrorHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
