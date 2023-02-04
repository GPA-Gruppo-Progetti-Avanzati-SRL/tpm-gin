package middleware

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
	"net/http"
	"reflect"
)

type AppError interface {
	GetCode() int
	GetMessage() string
	Error() string
	Marshal(ct string) ([]byte, error)
	Sanitized() AppError
}

type appError struct {
	Code        int    `json:"code,omitempty" yaml:"code,omitempty" mapstructure:"code,omitempty"`
	Ambit       string `json:"ambit,omitempty" yaml:"ambit,omitempty" mapstructure:"ambit,omitempty"`
	Text        string `json:"text,omitempty" yaml:"text,omitempty" mapstructure:"text,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
}

func (ae appError) Error() string {
	return ae.Text
}

func (ae appError) GetCode() int {
	return ae.Code
}

func (ae appError) GetMessage() string {
	return ae.Text
}

func (ae appError) Marshal(ct string) ([]byte, error) {

	if ct == "application/json" {
		b, err := json.Marshal(ae)
		return b, err
	}

	return nil, errors.New("app error cannot marshal to " + ct)
}

func (ae appError) Sanitized() AppError {

	nae := &appError{
		Code: ae.Code,
		Text: ae.Text,
	}

	return nae
}

func NewAppError(c int, m string) AppError {
	ae := &appError{Code: c, Text: m}
	return ae
}

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

				c.AbortWithStatusJSON(ae.GetCode(), ae)
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
		parsedError = &appError{
			Code: http.StatusInternalServerError,
			Text: "Internal Server Error",
		}

		return parsedError
	}

	switch v := gerr.Err.(type) {
	case AppError:
		parsedError = v
	default:
		parsedError = &appError{
			Code: http.StatusInternalServerError,
			Text: "Internal Server Error",
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
