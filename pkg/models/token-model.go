package models

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type TokenModel interface {
	Create(context.Context, uint64) error
	GetByAtUUID(context.Context, uint64, string) (uint64, error)
	DeleteByUUID(context.Context, string, string) error
}

type TokenDetails struct {
	AToken    string
	RToken    string
	AtUUID    string
	RtUUID    string
	AtExpires int64
	RtExpires int64
}

func (td *TokenDetails) Create(ctx context.Context, userID uint64) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := rds.Set(ctx, td.AtUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := rds.Set(ctx, td.RtUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (td *TokenDetails) GetByAtUUID(ctx context.Context, usrID uint64, atUUID string) (uint64, error) {
	userid, err := rds.Get(ctx, atUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	if usrID != userID {
		return 0, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized))
	}
	return userID, nil
}

func (td *TokenDetails) DeleteByUUID(ctx context.Context, atUUID, rtUUID string) error {
	ok, err := rds.Del(ctx, atUUID).Result()
	if ok != 1 {
		return err
	}
	ok, err = rds.Del(ctx, rtUUID).Result()
	if ok != 1 {
		return err
	}
	return nil
}
