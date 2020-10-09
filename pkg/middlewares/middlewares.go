package middlewares

import (
	"errors"
	"strconv"

	"github.com/user-service/pkg/auth"
	. "github.com/user-service/pkg/models"
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

const (
	corsAllowHeaders     = "Authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		next(ctx)
	}
}

func AUTH(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token, err := auth.ExtractAtMetadata(ctx)
		if err != nil {
			ERROR(ctx, fasthttp.StatusUnauthorized, err)
			return
		}

		var tokenDetails TokenDetails
		_, err = tokenDetails.GetByAtUUID(ctx, token.UserID, token.AtUUID)
		if err != nil {
			ERROR(ctx, fasthttp.StatusUnauthorized, err)
			return
		}
		ctx.SetUserValue(UserID, token.UserID)
		ctx.SetUserValue(AtUUID, token.AtUUID)

		next(ctx)
	}
}

func TRUTH(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		userId, err := strconv.ParseUint(ctx.UserValue("userId").(string), 10, 64)
		if err != nil {
			ERROR(ctx, fasthttp.StatusForbidden, errors.New(fasthttp.StatusMessage(fasthttp.StatusForbidden)))
			return
		}

		userID := ctx.UserValue(UserID)
		if userId != userID {
			ERROR(ctx, fasthttp.StatusForbidden, errors.New(fasthttp.StatusMessage(fasthttp.StatusForbidden)))
			return
		}

		next(ctx)
	}
}
