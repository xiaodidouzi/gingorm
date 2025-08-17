package tasks

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"context"
	"strconv"
	"strings"
	"time"
)

func SyncLikesToDB() {
	ctx := context.Background()
	iter := global.RedisDB.Scan(ctx, 0, "article:*likes", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		parts := strings.Split(key, ":")
		idStr := strings.TrimSuffix(parts[1], "likes")
		articleID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		likesStr, err := global.RedisDB.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		likes, _ := strconv.ParseUint(likesStr, 10, 64)
		global.DB.Model(&models.Article{}).Where("id = ?", articleID).Update("likes", likes)
	}
}
func StartLikesSyncTicker() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			SyncLikesToDB()
		}
	}()
}
