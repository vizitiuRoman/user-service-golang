package models

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
		create table if not exists users (
			id serial primary key,
			email varchar(255) NOT NULL UNIQUE,
			password varchar(255) NOT NULL,
			age varchar(255) NOT NULL,
			first_name varchar(255) NOT NULL,
			last_name varchar(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`

var (
	db  *sqlx.DB
	rds *redis.Client
)

func InitDatabase() error {
	DBSpec := fmt.Sprintf(
		"user=%s dbname=%s password=%s port=%s host=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
	)
	database, err := sqlx.Connect(os.Getenv("DB_DRIVER"), DBSpec)
	if err != nil {
		return err
	}
	db = database
	db.MustExec(schema)
	return nil
}

func InitRedis() error {
	host, port := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
	if len(host) == 0 {
		return errors.New("REDIS_HOST env does not exist")
	} else if len(port) == 0 {
		return errors.New("REDIS_PORT env does not exist")
	}
	rds = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})
	_, err := rds.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}
