package service

import (
	"context"
	"fmt"
	"gingorm/kafka"
	"gingorm/service/dto"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

type LikeService struct {
	DB            *gorm.DB
	Rdb           *redis.Client
	KafkaProducer *kafka.LikeProducer
}

func NewLikeService(db *gorm.DB, rdb *redis.Client, producer *kafka.LikeProducer) *LikeService {
	return &LikeService{DB: db, Rdb: rdb, KafkaProducer: producer}
}

func (s *LikeService) Like(ctx context.Context, userID uint, req dto.LikeRequest) (string, error) {
	key := fmt.Sprintf("like:%s:%d", req.TargetType, req.TargetID)
	field := fmt.Sprintf("%d", userID)

	exists, err := s.Rdb.HExists(ctx, key, field).Result()
	if err != nil {
		return "", err
	}

	var action string
	if exists {
		if err := s.Rdb.HDel(ctx, key, field).Err(); err != nil {
			return "", err
		}
		action = "unlike"
	} else {
		if err := s.Rdb.HSet(ctx, key, field, 1).Err(); err != nil {
			return "", err
		}
		action = "like"
	}

	msg := dto.LikeMessage{
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		Action:     action,
	}

	go s.KafkaProducer.SendLike(msg)

	return map[string]string{"like": "点赞成功", "unlike": "取消点赞成功"}[action], nil
}
