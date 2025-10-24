package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")

	if dbStr == "" {
		dbStr = "0"
	}

	db, err := strconv.Atoi(dbStr)

	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB value: %w", err)
	}

	rprotocol := os.Getenv("REDIS_PROTOCOL")

	protocol, err := strconv.Atoi(rprotocol)

	if err != nil {
		protocol = 2
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: protocol,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

func CloseRedis(rdb *redis.Client) {
	if rdb != nil {
		rdb.Close()
	}
}
