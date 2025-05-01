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

	// Goroutine to delete logs in scheduler
	go func() {
		for {
			now := time.Now()

			// Set next 4 AM
			next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
			if now.After(next) {
				// Already past today's 4 AM, schedule for tomorrow
				next = next.Add(24 * time.Hour)
			}
			duration := next.Sub(now)

			log.Printf("[Cleaner] Sleeping for %v until next 4AM...", duration)

			time.Sleep(duration)
			fmt.Println("[Cleaner] Running scheduled log cleanup...")

			var daysBefore int = 60

			settings, errs := controller.GetCoreSystemSettings()

			if errs != nil {
				log.Print("[Error] in fetching system settings")
			}

			if settings != nil {
				daysBefore = settings.AutoLogDeleteDays
			}

			from := time.Now().AddDate(0, 0, -daysBefore)
			to := time.Now()

			if err := controller.DeleteLogsScheduled(from, to); err != nil {
				log.Println("[Cleaner] Log delete error:", err)
			} else {
				log.Println("[Cleaner] Log cleanup completed.")
			}
		}
	}()

	// Normal worker to save logs in db
	for {
		logEntrys := []*request.LogEntry{}
		result, err := rdb.LPop(wrokerctx, "logs").Result()

		if err == redis.Nil {
			time.Sleep(100 * time.Millisecond)
			continue
		} else if err != nil {
			log.Println("Redis error:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// fmt.Printf("Received from Redis: %s\n", result)

		if err := json.Unmarshal([]byte(result), &logEntrys); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		if err := controller.SaveLogsBulkToDB(logEntrys); err != nil {
			log.Println("DB insert error:", err)
			continue
		}
	}
}
