package controllers

import "go.uber.org/zap"

type Controller struct {
	logger *zap.SugaredLogger
}

func NewController(logger *zap.SugaredLogger) *Controller {
	return &Controller{logger}
}
