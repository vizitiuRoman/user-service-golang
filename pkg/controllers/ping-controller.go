package controllers

import (
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

func (c *Controller) Ping(ctx *fasthttp.RequestCtx) {
	c.logger.Info("Ping")
	JSON(ctx, fasthttp.StatusOK, "User service OK")
}
