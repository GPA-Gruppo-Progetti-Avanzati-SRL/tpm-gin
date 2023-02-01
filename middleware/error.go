package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AppError interface {
	GetCode() int
	GetMessage() string
	Error() string
}

type appError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ae appError) Error() string {
	return ae.Message
}

func (ae appError) GetCode() int {
	return ae.Code
}

func (ae appError) GetMessage() string {
	return ae.Message
}

func NewAppError(c int, m string) AppError {
	ae := &appError{Code: c, Message: m}
	return ae
}

type ErrorHandler struct {
	config *ErrorHandlerConfig
}

func NewErrorHandler(cfg interface{}) MiddlewareHandler {
	var tcfg *ErrorHandlerConfig
	var ok bool
	if tcfg, ok = cfg.(*ErrorHandlerConfig); !ok {
		tcfg = &DefaultErrorHandlerConfig
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

			/*
			 * Don't know if should check also on c.Request.ContextParams().Err()
			 */
			if len(c.Errors) > 0 {
				for _, e := range c.Errors {
					log.Error().Str("middleware", "error").Msg(e.Error())
				}

				ae := getAppError(c.Errors[0])
				log.Error().Interface("appError", ae).Send()
				if h.config.WithInfo {
					c.AbortWithStatusJSON(ae.GetCode(), ae)
				} else {
					c.AbortWithStatus(ae.GetCode())
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
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}

		return parsedError
	}

	switch v := gerr.Err.(type) {
	case AppError:
		parsedError = v
	default:
		parsedError = &appError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return parsedError
}
