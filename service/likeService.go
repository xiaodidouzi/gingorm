package service

import (
	"context"
	"fmt"
	"gingorm/kafka"
	"gingorm/service/dto"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

type LikeService struct {
	DB            *gorm.DB
	Rdb           *redis.Client
	KafkaProducer *kafka.Producer
}

func NewLikeService(db *gorm.DB, rdb *redis.Client, kafka *kafka.Producer) *LikeService {
	return &LikeService{DB: db, Rdb: rdb, KafkaProducer: kafka}
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
		// 取消点赞
		if err := s.Rdb.HDel(ctx, key, field).Err(); err != nil {
			return "", err
		}
		action = "unlike"
	} else {
		// 点赞
		if err := s.Rdb.HSet(ctx, key, field, 1).Err(); err != nil {
			return "", err
		}
		action = "like"
	}

	// 发送 Kafka 消息
	msg := dto.LikeMessage{
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		Action:     action,
	}
	kafkaCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = s.KafkaProducer.SendMessage(kafkaCtx, msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	return map[string]string{"like": "点赞成功", "unlike": "取消点赞成功"}[action], nil
}
