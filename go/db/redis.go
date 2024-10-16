package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	ctx    = context.Background()
	client *redis.Client
)

// ConnectRedis initializes the Redis client using Upstash credentials
func ConnectRedis() {
	var err error
	opt, err := redis.ParseURL("rediss://default:AYzCAAIncDFmYjM0ZmJmMDg5OTk0Nzg4OWM3MzIzMTdkOTc3ZmRhNnAxMzYwMzQ@moving-bass-36034.upstash.io:6379")
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	client = redis.NewClient(opt)
}

// setRedis sets a key-value pair in Redis
func SetRedis(key, value string) error {
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
		return err
	}
	return nil
}
func CloseRedis() {
	if err := client.Close(); err != nil {
		log.Printf("Failed to close Redis client: %v", err)
	}
}
