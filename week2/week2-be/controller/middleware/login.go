package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

func (*LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/account/signup" || path == "/account/login" {
			return
		}
		s := sessions.Default(ctx)
		if s.Get("account_id") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
