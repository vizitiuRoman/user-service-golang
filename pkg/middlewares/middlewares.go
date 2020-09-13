package middlewares

import (
	"errors"

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
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		next(ctx)
	}
}

func AUTH(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token, err := auth.ExtractTokenMetadata(ctx)
		if err != nil {
			ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
			return
		}

		tokenDetails := TokenDetails{AccessUUID: token.AccessUUID, RefreshUUID: token.RefreshUUID}
		_, err = tokenDetails.GetByUUID(ctx, token.UserID)
		if err != nil {
			ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
			return
		}

		ctx.SetUserValue(UserID, token.UserID)
		ctx.SetUserValue(AccessUUID, token.AccessUUID)
		ctx.SetUserValue(RefreshUUID, token.RefreshUUID)

		next(ctx)
	}
}
