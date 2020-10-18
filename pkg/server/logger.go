package server

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.SugaredLogger {
	var config zap.Config
	if os.Getenv("LOG_LEVEL") == "prod" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger.Sugar()
}
