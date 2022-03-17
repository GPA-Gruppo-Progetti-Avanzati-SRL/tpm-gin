package middleware

const (
	MiddlewareErrorId                 = "gin-mw-error"
	MiddlewareErrorKind               = "mw-kind-error"
	MiddlewareErrorDefaultDiscoleInfo = true
)

/*
 * ErrorHandlerConfig
 */
type ErrorHandlerConfig struct {
	DiscloseErrorInfo bool `yaml:"disclose-error-info"  mapstructure:"disclose-error-info"`
}

var DefaultErrorHandlerConfig = ErrorHandlerConfig{
	DiscloseErrorInfo: MiddlewareErrorDefaultDiscoleInfo,
}

func (h *ErrorHandlerConfig) GetKind() string {
	return MiddlewareErrorKind
}

type ErrorHandlerConfigOption func(*ErrorHandlerConfig)
type ErrorHandlerConfigBuilder struct {
	opts []ErrorHandlerConfigOption
}

func (cb *ErrorHandlerConfigBuilder) WithErrorEnabled(enabled bool) *ErrorHandlerConfigBuilder {

	f := func(c *ErrorHandlerConfig) {
		c.DiscloseErrorInfo = enabled
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
