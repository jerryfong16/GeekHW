package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct{}

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

var JWTSecret = []byte("demo-app-hw")

type AccountJWTClaims struct {
	jwt.RegisteredClaims
	AccountId int64
	UserAgent string
}

type LoginJWTMiddlewareBuilder struct{}

func (*LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/account/signup" || path == "/account/login" {
			return
		}
		authorizationHeaderValue := ctx.GetHeader("Authorization")
		if authorizationHeaderValue == "" {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		segments := strings.Split(authorizationHeaderValue, " ")
		if len(segments) != 2 {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		tokenStr := segments[1]
		var claims AccountJWTClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expiredTime := claims.ExpiresAt
		if expiredTime.Sub(time.Now()) < time.Minute*2 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			tokenStr, err = token.SignedString(JWTSecret)
			ctx.Header("x-jwt", tokenStr)
			if err != nil {
			}
		}
		ctx.Set("account", claims)
	}
}
