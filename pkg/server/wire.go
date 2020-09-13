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

	router.GET("/api/v1/home", Home)

	// Auth
	router.POST("/api/v1/auth/login", Login)
	router.POST("/api/v1/auth/register", Register)
	router.GET("/api/v1/auth/logout/{id}", middlewares.AUTH(Logout))

	router.GET("/api/v1/token", middlewares.AUTH(Token))

	return router
}
