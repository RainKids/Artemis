package auth

import (
	"blog/pkg/exception"
	"blog/pkg/jwt"
	"blog/pkg/tools/network"
	"blog/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func Next() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		refer := ctx.Request.Header.Get("Referer")
		reqURI := ctx.Request.RequestURI
		// 验证Referer
		if network.GetReferDomain(refer) != network.GetReferDomain(reqURI) {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": "Referer验证失败"})
			ctx.Abort()
			return
		}
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusAccepted)
		}
		ctx.Next()
	}
}

func AuthenticateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			response.FailedResponse(ctx, exception.NewUnauthorized("Authorization error, %s", errors.New("No Authorization Token")))
			return
		}
		token := strings.Fields(authorization)[1]
		claim, err := jwt.ParseToken(token)
		if err != nil {
			if err, ok := err.(*gojwt.ValidationError); ok {
				if err.Errors&gojwt.ValidationErrorMalformed != 0 {
					response.FailedResponse(ctx, exception.NewUnauthorized("Authorization error, %s", err))
					ctx.Abort()
					return
				}
				if err.Errors&gojwt.ValidationErrorExpired != 0 {
					response.FailedResponse(ctx, exception.NewAccessTokenExpired("Authorization error, %s", err))
					ctx.Abort()
					return
				}
				if err.Errors&gojwt.ValidationErrorNotValidYet != 0 {
					response.FailedResponse(ctx, exception.NewAccessTokenIllegal("Authorization error, %s", err))
					ctx.Abort()
					return
				}
			}
		}
		if claim != nil {
			ctx.Set("username", claim.Username)
			ctx.Set("userId", claim.UserId)
			ctx.Next()
			return
		}
		response.FailedResponse(ctx, exception.NewUnauthorized("Authorization Failed, %s", err))
		ctx.Abort()
		return
	}
}
