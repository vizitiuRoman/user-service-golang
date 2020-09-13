package controllers

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/user-service/pkg/auth"
	. "github.com/user-service/pkg/models"
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

type UserWithToken struct {
	*User
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

	JSON(ctx, fasthttp.StatusOK, UserWithToken{foundUser, token})
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

	JSON(ctx, fasthttp.StatusCreated, UserWithToken{createdUser, token})
}

func Logout(ctx *fasthttp.RequestCtx) {
	id, err := strconv.ParseUint(ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		ERROR(ctx, fasthttp.StatusBadRequest, err)
		return
	}

	userID := ctx.UserValue(UserID)
	aUUID := ctx.UserValue(AccessUUID).(string)
	rUUID := ctx.UserValue(RefreshUUID).(string)

	if id != userID {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}

	var token TokenDetails
	err = token.DeleteByUUID(ctx, aUUID, rUUID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}
	JSON(ctx, fasthttp.StatusOK, fasthttp.StatusMessage(fasthttp.StatusOK))
}
