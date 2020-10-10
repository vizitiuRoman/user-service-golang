package auth

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/user-service/pkg/models"
	"github.com/valyala/fasthttp"
)

const (
	userID     = uint64(1)
	mockBearer = "rewiorweir2349"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env.test"))
	if err != nil {
		log.Fatalf("Cannot load env error: %v", err)
	}
	err = models.InitRedis()
	if err != nil {
		log.Fatalf("Cannot init redis error: %v", err)
	}
	m.Run()
}

func TestCreateToken(t *testing.T) {
	ctx := context.Background()
	createdToken, err := CreateToken(ctx, userID)
	if err != nil {
		t.Errorf("Cannot create token error: %v", err)
	}
	assert.NotEmpty(t, createdToken.AToken)
	assert.NotEmpty(t, createdToken.RToken)
}

func TestExtractToken(t *testing.T) {
	var ctx fasthttp.RequestCtx
	ctx.Request.Header.Set("Authorization", "Bearer "+mockBearer)

	token := extractToken(&ctx)

	assert.NotEmpty(t, token)
	assert.Equal(t, token, mockBearer)
}

func TestExtractAtMetadata(t *testing.T) {
	ctx := context.Background()
	createdToken, err := CreateToken(ctx, userID)
	if err != nil {
		t.Errorf("Cannot create token error: %v", err)
	}
	assert.NotEmpty(t, createdToken.AToken)
	assert.NotEmpty(t, createdToken.RToken)

	var ctxt fasthttp.RequestCtx
	ctxt.Request.Header.Set("Authorization", "Bearer "+createdToken.AToken)

	atDetails, err := ExtractAtMetadata(&ctxt)
	if err != nil {
		t.Errorf("Cannot extract access token error: %v", err)
	}

	assert.NotEmpty(t, atDetails.AtUUID)
	assert.Equal(t, atDetails.UserID, userID)
}

func TestExtractRtMetadata(t *testing.T) {
	ctx := context.Background()
	createdToken, err := CreateToken(ctx, userID)
	if err != nil {
		t.Errorf("Cannot create token error: %v", err)
	}
	assert.NotEmpty(t, createdToken.AToken)
	assert.NotEmpty(t, createdToken.RToken)

	atDetails, err := ExtractRtMetadata(createdToken.RToken)
	if err != nil {
		t.Errorf("Cannot extract refresh token error: %v", err)
	}

	assert.NotEmpty(t, atDetails.RtUUID)
	assert.Equal(t, atDetails.UserID, userID)
}
