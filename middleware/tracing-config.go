package middleware

const (
	MiddlewareTracingId              = "gin-mw-tracing"
	MiddlewareTracingKind            = "mw-kind-tracing"
	MiddlewareTracingDefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-"
	MiddlewareTracingDefaultSpanTag  = "error.id"
	MiddlewareTracingDefaultHeader   = "x-errid"
)

/*
 * TracingHandlerConfig
 */
type TracingHandlerConfig struct {
	Alphabet string
	SpanTag  string
	Header   string
}

var DefaultTracingHandlerConfig = TracingHandlerConfig{
	Alphabet: MiddlewareTracingDefaultAlphabet,
	SpanTag:  MiddlewareTracingDefaultSpanTag,
	Header:   MiddlewareTracingDefaultHeader,
}

func (h *TracingHandlerConfig) GetKind() string {
	return MiddlewareTracingKind
}

//    WithErrorDisclosureEnabled(bool)         // Enables/Disables error disclosure to the client
//                                             // if enabled the http error description is propagated to the client
//                                             // if disabled a response Header, configured with WithErrorDisclosureHeader is returned
//                                             // to the client with an errorid and the error is injected in an opentracing span having
//                                             // the same id as tag
//    WithErrorDisclosureSpanTag(string)       // span tag for the error  (defaults to "error.id)
//    WithErrorDisclosureHeader(string)        // error id header (defaults to "x-errid")
//    WithAlphabet(string)                     // alphabet  to generate the error id
type TracingHandlerConfigOption func(*TracingHandlerConfig)
type TracingHandlerConfigBuilder struct {
	opts []TracingHandlerConfigOption
}

func (cb *TracingHandlerConfigBuilder) WithAlphabet(alphabet string) *TracingHandlerConfigBuilder {
	f := func(c *TracingHandlerConfig) {
		c.Alphabet = alphabet
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *TracingHandlerConfigBuilder) WithSpanTag(s string) *TracingHandlerConfigBuilder {
	f := func(c *TracingHandlerConfig) {
		c.SpanTag = s
	}

	cb.opts = append(cb.opts, f)
	return cb
}
func (cb *TracingHandlerConfigBuilder) WithHeader(h string) *TracingHandlerConfigBuilder {
	f := func(c *TracingHandlerConfig) {
		c.Header = h
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *TracingHandlerConfigBuilder) Build() *TracingHandlerConfig {
	c := DefaultTracingHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
