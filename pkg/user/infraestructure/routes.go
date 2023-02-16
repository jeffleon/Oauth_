package infraestructure

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffleon/oauth-microservice/cmd/middleware"
)

type UserRoutes struct {
	handler UserHandler
}

func (ro *UserRoutes) PublicRoutes(public *gin.RouterGroup) {
	public.POST("/user/signup", ro.handler.SignUp)
	public.POST("/user/signin", ro.handler.SignIn)
}

func (ro *UserRoutes) SecureRoutes(secure *gin.RouterGroup) {
	secure.Use(middleware.TokenMiddleware(ro.handler.Service.VerifyToken, ro.handler.Service.VerifyBlackList))
	secure.GET("/user", ro.handler.GetUser)
	secure.POST("/user/logout", ro.handler.Logout)
	secure.POST("/user/refreshToken", ro.handler.RefreshAccessToken)
	// secure.PATCH("/user", ro.handler.UpdateUser)
	// secure.DELETE("/user", ro.handler.DeleteUser)

}

func NewRoutes(handler UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}
