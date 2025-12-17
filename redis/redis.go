package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client

func NewRedisClient() *redis.Client {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	log.Println("Redis connected")
	return RDB
}

// Increment trending product
func IncrementTrending(productId string) error {
	return RDB.ZIncrBy(Ctx, "trending_products", 1, productId).Err()
}

// Get Top N Trending
func GetTopTrending(n int) ([]redis.Z, error) {
	return RDB.ZRevRangeWithScores(Ctx, "trending_products", 0, int64(n-1)).Result()
}
