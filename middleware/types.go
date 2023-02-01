package middleware

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareHandler interface {
	GetKind() string
	HandleFunc() gin.HandlerFunc
}

/*
 * Package Configuration defaults

func GetConfigDefaults(contextPath string) []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{strings.Join([]string{contextPath, ErrorHandlerId, "disclose-error-info"}, "."), ErrorHandlerDefaultWithInfo, "error is in clear"},
		{strings.Join([]string{contextPath, TracingHandlerId, "alphabet"}, "."), TracingHandlerDefaultAlphabet, "alphabet"},
		{strings.Join([]string{contextPath, TracingHandlerId, "spantag"}, "."), TracingHandlerDefaultSpanTag, "spantag"},
		{strings.Join([]string{contextPath, TracingHandlerId, "header"}, "."), TracingHandlerDefaultHeader, "header"},
	}
}
*/
