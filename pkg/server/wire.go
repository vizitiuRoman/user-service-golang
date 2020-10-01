package server

import (
	"os"

	. "github.com/fasthttp/router"
	. "github.com/user-service/pkg/controllers"
	"github.com/user-service/pkg/middlewares"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() error {
	var config zap.Config
	if os.Getenv("LOG_LEVEL") == "prod" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	zap.ReplaceGlobals(logger)
	if err != nil {
		return err
	}
	return nil
}

func initControllers() *Router {
	router := New()

	// Ping
	router.GET("/api/ping", Ping)

	// Auth
	router.POST("/api/v1/auth/login", Login)
	router.POST("/api/v1/auth/register", Register)
	router.GET("/api/v1/auth/refresh", middlewares.AUTH(RefreshToken))
	router.GET("/api/v1/auth/logout/{userId}", middlewares.AUTH(middlewares.TRUTH(Logout)))

	// User
	router.GET("/api/v1/users", GetUsers)
	router.GET("/api/v1/users/{userId}", middlewares.AUTH(middlewares.TRUTH(GetUser)))
	router.POST("/api/v1/users", middlewares.AUTH(UpdateUser))
	router.DELETE("/api/v1/users/{userId}", middlewares.AUTH(middlewares.TRUTH(DeleteUser)))

	return router
}
