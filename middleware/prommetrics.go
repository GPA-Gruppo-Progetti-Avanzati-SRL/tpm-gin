package middleware

import (
	"github.com/gin-gonic/gin"
)

type PromHttpMetricsHandler struct {
	config *PromHttpMetricsHandlerConfig
}

func NewPromHttpMetricsHandler(cfg interface{}) MiddlewareHandler {
	var tcfg *PromHttpMetricsHandlerConfig
	var ok bool
	if tcfg, ok = cfg.(*PromHttpMetricsHandlerConfig); !ok {
		tcfg = &DefaultPromHttpMetricsHandlerConfig
	}

	return &PromHttpMetricsHandler{
		config: tcfg,
	}
}

func (h *PromHttpMetricsHandler) GetKind() string {
	return MiddlewareMetricsKind
}

func (m *PromHttpMetricsHandler) HandleFunc() gin.HandlerFunc {

	// TODO: at the moment it is simply an empty handler..
	return func(c *gin.Context) {
		if nil != c {
			c.Next()
		}
	}

}
