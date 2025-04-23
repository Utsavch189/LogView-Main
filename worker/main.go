package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Utsavch189/logview/internal/controller"
	"github.com/Utsavch189/logview/internal/models/request"
	"github.com/go-redis/redis/v8"
)

var wrokerctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	for {
		logEntry := &request.LogEntry{}
		result, err := rdb.LPop(wrokerctx, "logs").Result()

		if err == redis.Nil {
			time.Sleep(100 * time.Millisecond)
			continue
		} else if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		if err := json.Unmarshal([]byte(result), logEntry); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		fmt.Printf("[LOG][%s] %s - %s - %s\n", logEntry.Level, logEntry.Logger, logEntry.Message, logEntry.Exception)

		if err := controller.SaveLogToDB(logEntry); err != nil {
			log.Println("DB insert error:", err)
		} else {
			println("Db insert complete")
		}
	}
}
