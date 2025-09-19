package redisclient

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init() {
	// Load .env if exists
	_ = godotenv.Load()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	addr := fmt.Sprintf("%s:%s", host, port)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if _, err := Rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("⚠️ Redis unavailable: %v", err)
		Rdb = nil
	} else {
		log.Println("✅ Redis connected successfully")
	}
}
