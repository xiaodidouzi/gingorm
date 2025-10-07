package config

import (
	"context"
	"gingorm/global"
	"github.com/go-redis/redis/v8"
	"log"
)

func InitRedis() {
	addr := AppConfig.Redis.Addr
	db := AppConfig.Redis.DB
	password := AppConfig.Redis.Password
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})
	ctx := context.Background()
	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis,get error:%v", err)
	}
	global.RedisDB = RedisClient
}
