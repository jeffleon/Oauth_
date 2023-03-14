package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffleon/oauth-microservice/pkg/health"
	"github.com/jeffleon/oauth-microservice/pkg/swagger"
	users "github.com/jeffleon/oauth-microservice/pkg/user/infraestructure"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(routes RoutesGroup) Router {
	route := gin.Default()
	health := route.Group("api/OAuth/v1")
	group := route.Group("/")
	public := route.Group("api/OAuth/v1/public")
	secure := route.Group("api/OAuth/v1/secure")
	routes.User.PublicRoutes(public)
	routes.User.SecureRoutes(secure)
	routes.Health.RegisterRoutes(health, group)
	routes.Swagger.RegisterRoutes(public)
	return route
}

type RoutesGroup struct {
	User    *users.UserRoutes
	Health  *health.Routes
	Swagger *swagger.Routes
}
