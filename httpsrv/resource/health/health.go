package health

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-gin/httpsrv"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func init() {
	log.Info().Msg("health resources init function invoked")
	ra := httpsrv.GetApp()
	ra.RegisterGFactory(registerHealthEndpoints)
}

func registerHealthEndpoints(ctx httpsrv.ServerContext) []httpsrv.G {

	gs := make([]httpsrv.G, 0, 2)

	gs = append(gs, httpsrv.G{
		Name:    "Liveness endpoint",
		AbsPath: true,
		Path:    "health",
		Resources: []httpsrv.R{
			{
				Name:          "liveness",
				Path:          "liveness",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{func(c *gin.Context) { c.JSON(200, "OK") }},
			},
			{
				Name:          "readiness",
				Path:          "readiness",
				Method:        http.MethodGet,
				RouteHandlers: []httpsrv.H{func(c *gin.Context) { c.JSON(200, "OK") }},
			},
		},
	})

	return gs
}
