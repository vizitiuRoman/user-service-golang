package controllers

import (
	"github.com/user-service/pkg/auth"
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

func Home(ctx *fasthttp.RequestCtx) {
	token, err := auth.CreateToken(ctx, 10)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}
	JSON(ctx, fasthttp.StatusOK, token)
}

func Token(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue(UserID)
	JSON(ctx, fasthttp.StatusOK, userID)
}
