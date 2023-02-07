package middleware

import (
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
	"net/http"
	"reflect"
)

type ErrorHandler struct {
	config *ErrorHandlerConfig
}

func NewErrorHandler(cfg interface{}) MiddlewareHandler {
	var tcfg *ErrorHandlerConfig
	var ok bool

	if cfg == nil || reflect.ValueOf(cfg).IsNil() {
		tcfg = &DefaultErrorHandlerConfig
	} else {
		if tcfg, ok = cfg.(*ErrorHandlerConfig); !ok {
			tcfg = &DefaultErrorHandlerConfig
		}
	}

	return &ErrorHandler{
		config: tcfg,
	}
}

func (h *ErrorHandler) GetKind() string {
	return ErrorHandlerKind
}

func (h *ErrorHandler) HandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if nil != c {
			c.Next()

			if len(c.Errors) > 0 {
				for _, e := range c.Errors {
					log.Error().Str("middleware", "error").Msg(e.Error())
				}

				var ae AppError
				ae = getAppError(c.Errors[0])
				if !h.config.WithCause {
					ae = ae.Sanitized()
				}

				span := opentracing.SpanFromContext(c.Request.Context())
				if nil != span {
					h.fail(c, span, ae.Error())
				}

				c.AbortWithStatusJSON(ae.GetStatusCode(), ae)
			} else if h.config.StatusCodeHandlingPolicy.Hightlight(c.Writer.Status()) {
				span := opentracing.SpanFromContext(c.Request.Context())
				if nil != span {
					// In this case headers have been written already.... c.Writer.Written() == true
					// limit to setting trace span error flag
					ext.Error.Set(span, true)
				}
			}
		}
	}
}

func getAppError(err error) AppError {
	var parsedError AppError

	gerr, ok1 := err.(*gin.Error)
	if !ok1 {
		parsedError = &AppErrorImpl{
			StatusCode: http.StatusInternalServerError,
			Text:       "Internal Server Error",
		}

		return parsedError
	}

	switch v := gerr.Err.(type) {
	case AppError:
		parsedError = v
	default:
		parsedError = &AppErrorImpl{
			StatusCode: http.StatusInternalServerError,
			Text:       "Internal Server Error",
			Message:    v.Error(),
		}
	}

	return parsedError
}

func (h *ErrorHandler) fail(c *gin.Context, span opentracing.Span, cause string) {
	ext.Error.Set(span, true)
	if cause != "" {
		span.SetTag("cause", cause)
	}

	// injecting error id and tagging span
	errid, err := gonanoid.Generate(h.config.Alphabet, 32)
	if nil != err {
		// in this case just dump error, we want error handling to be smooth
		// ignore
	} else {
		if nil != span {
			span.SetTag(h.config.SpanTag, errid)
			c.Header(h.config.Header, errid)
		}
	}
}
