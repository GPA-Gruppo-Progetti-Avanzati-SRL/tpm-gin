package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"reflect"
)

type TracingHandler struct {
	config *TracingHandlerConfig
}

// NewErrorHandler builds an Error Handler with the following options:

func NewTracingHandler(cfg interface{}) MiddlewareHandler {

	var tcfg *TracingHandlerConfig
	var ok bool

	if cfg == nil || reflect.ValueOf(cfg).IsNil() {
		tcfg = &DefaultTracingHandlerConfig
	} else {
		if tcfg, ok = cfg.(*TracingHandlerConfig); !ok {
			tcfg = &DefaultTracingHandlerConfig
		}
	}

	return &TracingHandler{
		config: tcfg,
	}
}

func (t *TracingHandler) GetKind() string {
	return TracingHandlerKind
}

func (t *TracingHandler) HandleFunc() gin.HandlerFunc {

	return func(c *gin.Context) {

		log.Trace().Str("requestPath", c.Request.RequestURI).Send()

		var span opentracing.Span
		parentSpanCtx, serr := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if nil != serr {
			span = opentracing.StartSpan(c.FullPath())
		} else {
			span = opentracing.StartSpan(c.FullPath(), opentracing.ChildOf(parentSpanCtx))
		}
		defer span.Finish()

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		if nil != c {
			c.Next()
		}

		statusCode := c.Writer.Status()
		if span != nil {
			span.SetTag("http.method", c.Request.Method)
			span.SetTag("http.status_code", statusCode)
		}

	}
}

//func (t *TracingHandler) fail(c *gin.Context, appErr AppError, span opentracing.Span) {
//
//	if nil != span {
//		ext.Error.Set(span, true)
//		span.SetTag("cause", appErr)
//		ext.HTTPStatusCode.Set(span, uint16(appErr.GetCode()))
//	}
//
//	// injecting error id and tagging span
//	errid, err := gonanoid.Generate(t.config.Alphabet, 32)
//	if nil != err { // in this case just dump error, we want error handling to be smooth
//		// ignore
//	} else {
//		if nil != span {
//			span.SetTag(t.config.SpanTag, errid)
//			c.Header(t.config.Header, errid)
//		}
//	}
//}
//
//func (t *TracingHandler) failWithContext(c *gin.Context, w http.ResponseWriter, appErr AppError) {
//	span := opentracing.SpanFromContext(c.Request.Context())
//	t.fail(c, appErr, span)
//}
