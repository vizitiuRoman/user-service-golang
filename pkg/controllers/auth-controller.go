package controllers

import (
	"encoding/json"
	"errors"

	"github.com/user-service/pkg/auth"
	. "github.com/user-service/pkg/models"
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

type authResponse struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func Login(ctx *fasthttp.RequestCtx) {
	var user User
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		ERROR(ctx, fasthttp.StatusBadRequest, err)
		return
	}

	password := user.Password
	foundUser, err := user.FindByEmail(user.Email)
	if err != nil {
		ERROR(ctx, fasthttp.StatusBadRequest, errors.New("Invalid email or password"))
		return
	}

	err = VerifyPassword(foundUser.Password, password)
	if err != nil {
		ERROR(ctx, fasthttp.StatusBadRequest, errors.New("Invalid email or password"))
		return
	}

	token, err := auth.CreateToken(ctx, foundUser.ID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}

	JSON(ctx, fasthttp.StatusOK, authResponse{foundUser.ID, foundUser.Email, token})
}

func Register(ctx *fasthttp.RequestCtx) {
	var user User
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		ERROR(ctx, fasthttp.StatusBadRequest, err)
		return
	}

	createdUser, err := user.Create()
	if err != nil {
		ERROR(ctx, fasthttp.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateToken(ctx, createdUser.ID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}

	JSON(ctx, fasthttp.StatusOK, authResponse{createdUser.ID, createdUser.Email, token})
}

func Logout(ctx *fasthttp.RequestCtx) {
	aUUID := ctx.UserValue(AccessUUID).(string)
	rUUID := ctx.UserValue(RefreshUUID).(string)

	var token TokenDetails
	err := token.DeleteByUUID(ctx, aUUID, rUUID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}
	JSON(ctx, fasthttp.StatusOK, fasthttp.StatusMessage(fasthttp.StatusOK))
}

func RefreshToken(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue(UserID).(uint64)
	aUUID := ctx.UserValue(AccessUUID).(string)
	rUUID := ctx.UserValue(RefreshUUID).(string)

	var oldToken TokenDetails
	err := oldToken.DeleteByUUID(ctx, aUUID, rUUID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}

	token, err := auth.CreateToken(ctx, userID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}

	var user User
	foundUser, err := user.FindByID(userID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}

	JSON(ctx, fasthttp.StatusOK, authResponse{foundUser.ID, foundUser.Email, token})
}
