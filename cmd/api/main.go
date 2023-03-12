package main

// @title     Oauth API
// @version   1.0

// @host      localhost:8080
// @BasePath  /api/Oauth/v1

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jeffleon/oauth-microservice/internal/config"
	"github.com/jeffleon/oauth-microservice/pkg/health"
	"github.com/jeffleon/oauth-microservice/pkg/router"
	"github.com/jeffleon/oauth-microservice/pkg/swagger"
	oauthService "github.com/jeffleon/oauth-microservice/pkg/user/aplication"
	oauthDomain "github.com/jeffleon/oauth-microservice/pkg/user/domain"
	oauthInfra "github.com/jeffleon/oauth-microservice/pkg/user/infraestructure"
	"github.com/jeffleon/oauth-microservice/pkg/user/infraestructure/repository"
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
		logrus.Infof("Error, cannot connect with rpc server %s", err)
		logrus.Infof("You can't send Emails")
	}

	kafkaProducer, err := InitKafkaProducer()

	if err != nil {
		panic(err)
	}
	logrus.Info("Kafka connected")
	defer kafkaProducer.Close()

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
	rpcRepo := repository.NewRPCRepository(clientRPC)
	kafkaRepo := repository.NewKafkaRepository(kafkaProducer)
	redisRepo := repository.NewRedisRepository(ctx, client)
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(tokenObj)
	userService := oauthService.NewUserService(userRepo, tokenRepo, redisRepo, rpcRepo, kafkaRepo)
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

func InitKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", config.Config.KafkaHost, config.Config.KafkaPort),
		"security.protocol": "SASL_SSL",
		"sasl.username":     config.Config.KafkaUsername,
		"sasl.password":     config.Config.KafkaPassword,
		"sasl.mechanism":    "PLAIN",
	})
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.RedisHost, config.Config.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
