package configs

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedisClient(addr, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully")
}
