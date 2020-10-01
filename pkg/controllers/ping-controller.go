package controllers

import (
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

func Ping(ctx *fasthttp.RequestCtx) {
	JSON(ctx, fasthttp.StatusOK, "User service OK")
}
