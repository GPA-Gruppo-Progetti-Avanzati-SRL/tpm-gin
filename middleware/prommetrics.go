package middleware

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/middleware/promutil"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"time"
)

type PromHttpMetricsHandler struct {
	config     *PromHttpMetricsHandlerConfig
	collectors []promutil.MetricInfo
}

func NewPromHttpMetricsHandler(cfg interface{}) MiddlewareHandler {
	var tcfg *PromHttpMetricsHandlerConfig
	var ok bool
	if tcfg, ok = cfg.(*PromHttpMetricsHandlerConfig); !ok {
		tcfg = &DefaultPromHttpMetricsHandlerConfig
	}

	if tcfg.Namespace == "" || tcfg.Subsystem == "" {
		tcfg = &DefaultMetricsConfig
	} else {
		if len(tcfg.Collectors) == 0 {
			tcfg.Collectors = DefaultMetricsConfig.Collectors
		}
	}

	collectors := make([]promutil.MetricInfo, 0)

	for _, mCfg := range tcfg.Collectors {
		if mc, err := promutil.NewCollector(tcfg.Namespace, tcfg.Subsystem, mCfg.Name, &mCfg); err != nil {
			log.Error().Err(err).Str("name", mCfg.Name).Msg("error creating metric")
		} else {
			collectors = append(collectors, promutil.MetricInfo{Type: mCfg.Type, Id: mCfg.Id, Name: mCfg.Name, Collector: mc, Labels: mCfg.Labels})
		}
	}

	return &PromHttpMetricsHandler{
		config:     tcfg,
		collectors: collectors,
	}
}

func (h *PromHttpMetricsHandler) GetKind() string {
	return MetricsHandlerKind
}

func (m *PromHttpMetricsHandler) HandleFunc() gin.HandlerFunc {

	return func(c *gin.Context) {

		beginOfMiddleware := time.Now()

		var sc = "500"
		ep := c.Request.URL.String()

		defer func(begin time.Time) {
			promutil.SetMetricValueById(m.collectors, "request_duration", time.Since(begin).Seconds(), prometheus.Labels{"endpoint": ep, "status_code": sc})
		}(beginOfMiddleware)

		if nil != c {
			c.Next()
		}

		sc = fmt.Sprintf("%d", c.Writer.Status())
		_ = promutil.SetMetricValueById(m.collectors, "requests", 1, prometheus.Labels{"endpoint": ep, "status_code": sc})
	}

}
