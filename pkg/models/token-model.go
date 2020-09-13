package models

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type TokenModel interface {
	Create(ctx context.Context, userID uint64) error
	GetByUUID(ctx context.Context, usrID uint64) (uint64, error)
	DeleteByUUID(ctx context.Context, aUUID, rUUID string) error
}

type TokenDetails struct {
	AccessToken string
	AccessUUID  string
	RefreshUUID string
	AtExpires   int64
	RtExpires   int64
}

func (t *TokenDetails) Create(ctx context.Context, userID uint64) error {
	at := time.Unix(t.AtExpires, 0)
	rt := time.Unix(t.RtExpires, 0)
	now := time.Now()

	errAccess := rds.Set(ctx, t.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := rds.Set(ctx, t.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (t *TokenDetails) GetByUUID(ctx context.Context, usrID uint64) (uint64, error) {
	userid, err := rds.Get(ctx, t.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	_, err = rds.Get(ctx, t.RefreshUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	if usrID != userID {
		return 0, errors.New(fasthttp.StatusMessage(fasthttp.StatusUnauthorized))
	}
	return userID, nil
}

func (t *TokenDetails) DeleteByUUID(ctx context.Context, aUUID, rUUID string) error {
	ok, _ := rds.Del(ctx, aUUID).Result()
	if ok != 1 {
		return errors.New("Error to delete access uuid")
	}
	ok, _ = rds.Del(ctx, rUUID).Result()
	if ok != 1 {
		return errors.New("Error to delete refresh uuid")
	}
	return nil
}
