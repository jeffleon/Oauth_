package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
)

type validateTokenService func(string) (*domain.Claims, bool)
type validateBalckListService func(string) bool

func TokenMiddleware(validateToken validateTokenService, validateBlackList validateBalckListService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")

		claims, validToken := validateToken(token)
		if !validToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		valid := validateBlackList(token)
		if !valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("user_id", claims.UserID)
		ctx.Set("authorized", claims.Authorized)

		ctx.Next()
	}
}
