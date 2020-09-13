package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/user-service/pkg/models"
	. "github.com/user-service/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"github.com/valyala/fasthttp"
)

type AccessDetails struct {
	AccessUUID  string
	RefreshUUID string
	UserID      uint64
}

func prepareToken(extractedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(extractedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return &jwt.Token{}, err
	}
	return token, nil
}

func extractToken(ctx *fasthttp.RequestCtx) string {
	bearerToken := ctx.Request.Header.Peek("Authorization")
	if len(strings.Split(string(bearerToken), " ")) == 2 {
		return strings.Split(string(bearerToken), " ")[1]
	}
	return ""
}

func CreateToken(ctx context.Context, userID uint64) (string, error) {
	accessUUID := uuid.NewV4().String()
	refreshUUID := uuid.NewV4().String()
	claims := jwt.MapClaims{}
	claims[UserID] = userID
	claims[AccessUUID] = accessUUID
	claims[RefreshUUID] = refreshUUID
	claims["exp"] = time.Now().Add(TokenExpires).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(
		[]byte(os.Getenv("API_SECRET")),
	)
	if err != nil {
		return "", err
	}

	tokenDetails := &TokenDetails{
		AccessToken: token,
		AccessUUID:  accessUUID,
		RefreshUUID: refreshUUID,
		AtExpires:   time.Now().Add(AtExpires).Unix(),
		RtExpires:   time.Now().Add(RtExpires).Unix(),
	}
	err = tokenDetails.Create(ctx, userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractTokenMetadata(ctx *fasthttp.RequestCtx) (*AccessDetails, error) {
	extractedToken := extractToken(ctx)
	if extractedToken == "" {
		return &AccessDetails{}, errors.New("Cannot extract token")
	}

	token, err := prepareToken(extractedToken)
	if err != nil {
		return &AccessDetails{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessUUID, ok := claims[AccessUUID].(string)
		if !ok {
			return &AccessDetails{}, errors.New("Cannot get access uuid")
		}
		refreshUUID, ok := claims[RefreshUUID].(string)
		if !ok {
			return &AccessDetails{}, errors.New("Cannot get refresh uuid")
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims[UserID]), 10, 32)
		if err != nil {
			return &AccessDetails{}, errors.New("Cannot get user id")
		}
		return &AccessDetails{
			AccessUUID:  accessUUID,
			RefreshUUID: refreshUUID,
			UserID:      userID,
		}, nil
	}
	return &AccessDetails{}, errors.New("ExtractTokenMetadata error")
}
