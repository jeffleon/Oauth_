package main

// @title     Oauth API
// @version   1.0

// @host      localhost:8080
// @BasePath  /api/Oauth/v1

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/jeffleon/oauth-microservice/internal/config"
	"github.com/jeffleon/oauth-microservice/pkg/health"
	"github.com/jeffleon/oauth-microservice/pkg/router"
	"github.com/jeffleon/oauth-microservice/pkg/swagger"
	oauthService "github.com/jeffleon/oauth-microservice/pkg/user/aplication"
	oauthDomain "github.com/jeffleon/oauth-microservice/pkg/user/domain"
	oauthInfra "github.com/jeffleon/oauth-microservice/pkg/user/infraestructure"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()
	db, err := InitDB()
	if err != nil {
		logrus.Fatalf("Error, Database cannot connect %s", err)
	}
	client := InitRedis()
	clientRPC, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", config.Config.RPCHost, config.Config.RPCPort))
	if err != nil {
		logrus.Fatalf("Error, cannot connect with rpc server %s", err)
	}
	tokenObj := oauthDomain.TokenObj{
		Predetermined: oauthDomain.TokenType{
			Secret:     config.Config.TokenSecret,
			Expiration: config.Config.TokenExp,
		},
		RefreshToken: oauthDomain.TokenType{
			Secret:     config.Config.TokenSecret,
			Expiration: config.Config.TokenRefreshExp,
		},
	}
	rpcRepo := oauthInfra.NewRPCRepository(clientRPC)
	redisRepo := oauthInfra.NewRedisRepository(ctx, client)
	userRepo := oauthInfra.NewUserRepository(db)
	tokenRepo := oauthInfra.NewTokenRepository(tokenObj)
	userService := oauthService.NewUserService(userRepo, tokenRepo, redisRepo, rpcRepo)
	userHandler := oauthInfra.UserHandler{Service: userService}
	userRoutes := oauthInfra.NewRoutes(userHandler)
	healthRoutes := health.NewHealthCheckRoutes()
	swaggerRoutes := swagger.NewSwaggerDocsRoutes()
	routes := router.RoutesGroup{
		User:    userRoutes,
		Health:  healthRoutes,
		Swagger: swaggerRoutes,
	}
	r := router.NewRouter(routes)
	logrus.Fatal(r.Run(fmt.Sprintf(":%s", config.Config.Port)))
}

func InitDB() (*gorm.DB, error) {
	dsn := config.Config.DBConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := oauthDomain.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
