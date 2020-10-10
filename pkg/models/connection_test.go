package models

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env.test"))
	if err != nil {
		log.Fatalf("Cannot load env error: %v", err)
	}
	m.Run()
}

func TestInitDatabase(t *testing.T) {
	err := InitDatabase()
	if err != nil {
		t.Fatalf("Cannot init database error: %v", err)
	}
	err = db.Ping()
	assert.Equal(t, err, nil)
}

func TestInitRedis(t *testing.T) {
	ctx := context.Background()

	err := InitRedis()
	if err != nil {
		t.Fatalf("Cannot init redis error: %v", err)
	}

	status := rds.Ping(ctx).Err()
	if status != nil {
		t.Fatalf("Cannot init redis error: %v", err)
	}

	valueErr := rds.Set(ctx, "value", "10", 0).Err()
	if valueErr != nil {
		t.Fatalf("Cannot set value error: %v", err)
	}

	value, err := rds.Get(ctx, "value").Result()
	if err != nil {
		t.Fatalf("Cannot get value error: %v", err)
	}
	assert.Equal(t, value, "10")
}
