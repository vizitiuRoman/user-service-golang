package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const userID = uint64(10)

func TestCreateToken(t *testing.T) {
	ctx := context.Background()
	td := TokenDetails{
		"qwe", "eqwe", "q41",
		"4324", 1, 1,
	}
	err := td.Create(ctx, userID)
	assert.Nil(t, err)
}

func TestGetByAtUUID(t *testing.T) {
	ctx := context.Background()
	td := TokenDetails{
		"10239", "349123", "1323u423",
		"1323u423", 1, 1,
	}
	err := td.Create(ctx, userID)
	if err != nil {
		t.Errorf("Cannot create token error: %v", err)
	}

	userId, err := td.GetByAtUUID(ctx, userID, td.AtUUID)
	if err != nil {
		t.Errorf("Cannot get userID by at uuid error: %v", err)
	}
	assert.Equal(t, userId, userID)
}

func TestDeleteByUUID(t *testing.T) {
	ctx := context.Background()
	td := TokenDetails{
		"102539", "349q123", "1323u4r23",
		"1323ru423", 1, 1,
	}
	err := td.Create(ctx, 100)
	if err != nil {
		t.Errorf("Cannot create token error: %v", err)
	}

	err = td.DeleteByUUID(ctx, td.AtUUID, td.RtUUID)
	if err != nil {
		t.Errorf("Cannot delete by uuid error: %v", err)
	}
}
