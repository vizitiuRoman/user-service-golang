package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/user-service/pkg/models"
	. "github.com/user-service/pkg/utils"
	"github.com/valyala/fasthttp"
)

func GetUsers(ctx *fasthttp.RequestCtx) {
	var user User
	users, err := user.FindAll()
	if err != nil {
		ERROR(ctx, fasthttp.StatusInternalServerError, err)
		return
	}
	JSON(ctx, fasthttp.StatusOK, users)
}

func GetUser(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue(UserID).(uint64)

	var user User
	foundUser, err := user.FindByID(userID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusNotFound, errors.New("Not found user"))
		return
	}
	JSON(ctx, fasthttp.StatusOK, foundUser)
}

func UpdateUser(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue(UserID).(uint64)

	var user User
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnprocessableEntity, err)
		return
	}

	if userID != user.ID {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		ERROR(ctx, fasthttp.StatusBadRequest, err)
		return
	}

	err = user.Update()
	if err != nil {
		ERROR(ctx, fasthttp.StatusInternalServerError, err)
		return
	}

	JSON(ctx, fasthttp.StatusOK, User{ID: userID, Email: user.Email})
}

func DeleteUser(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue(UserID).(uint64)
	atUUID := ctx.UserValue(AtUUID).(string)
	rtUUID := fmt.Sprintf("%s++%d", atUUID, userID)

	var token TokenDetails
	err := token.DeleteByUUID(ctx, atUUID, rtUUID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusUnauthorized, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized)))
		return
	}

	var user User
	err = user.DeleteByID(userID)
	if err != nil {
		ERROR(ctx, fasthttp.StatusInternalServerError, err)
		return
	}
	JSON(ctx, fasthttp.StatusOK, true)
}
