package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.alm.poste.it/go/configuration"
	"strings"
)

type MiddlewareHandler interface {
	GetKind() string
	HandleFunc() gin.HandlerFunc
}

/*
 * Package Configuration defaults
 */
func GetConfigDefaults(contextPath string) []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{strings.Join([]string{contextPath, MiddlewareErrorId, "disclose-error-info"}, "."), MiddlewareErrorDefaultDiscoleInfo, "error is in clear"},
		{strings.Join([]string{contextPath, MiddlewareTracingId, "alphabet"}, "."), MiddlewareTracingDefaultAlphabet, "alphabet"},
		{strings.Join([]string{contextPath, MiddlewareTracingId, "spantag"}, "."), MiddlewareTracingDefaultSpanTag, "spantag"},
		{strings.Join([]string{contextPath, MiddlewareTracingId, "header"}, "."), MiddlewareTracingDefaultHeader, "header"},
	}
}
