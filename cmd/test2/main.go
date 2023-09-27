package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"encoding/json"
	"github.com/go-playground/locales/lrc"
	"system/api/repository"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         "127.0.0.1:6379",
		Password:     "secret_redis",
		DB:           0,
		PoolSize:     1600,
		MaxRetries:   5,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	for i := 0; i < 100; i++ {
		redisClient.RPush(context.Background(), "lottery_win_message", strconv.Itoa(i)+"aaaa").Result()
	}
	lrc.New()

	// make cache with 10ms TTL and 5 max keys
	var userInfo repository.JavaUserInfo
	userInfo.UserId = "1"
	t, _ := json.Marshal(userInfo)

	err := json.Unmarshal(t, &userInfo)
	println(err == nil)
}
