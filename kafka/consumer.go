package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"

	"awesomeProject/global"
	"awesomeProject/models"

	"github.com/segmentio/kafka-go"
)

func StartLikeConsumer(brokers []string, topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID, // 消费组ID
		MinBytes: 10e3,    // 10KB
		MaxBytes: 10e6,    // 10MB
	})

	log.Println("Kafka Like Consumer started...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var msg LikeMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}
		err = global.DB.Transaction(func(tx *gorm.DB) error {
			like := models.Like{
				ArticleID: msg.ArticleID,
				UserID:    msg.UserID,
			}
			if err := global.DB.Create(&like).Error; err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					log.Printf("Duplicate like ignored: ArticleID=%d UserID=%d", msg.ArticleID, msg.UserID)
					return nil
				}
				return err
			}
			if err := tx.Model(&models.Article{}).Where("id=?", msg.ArticleID).
				UpdateColumn("likes", gorm.Expr("likes+?", 1)).Error; err != nil {
				return err
			}
			likeKey := fmt.Sprintf("like:%d:%d", msg.ArticleID, msg.UserID)
			if err := global.RedisDB.Set(context.Background(), likeKey, 1, 0).Err(); err != nil {
				log.Printf("Redis update failed: %v", err)
			}
			return nil
		})
		if err != nil {
			log.Printf("Transaction failed for ArticleID=%d UserID=%d: %v", msg.ArticleID, msg.UserID, err)
		} else {
			log.Printf("Like saved & Article.Likes updated: ArticleID=%d UserID=%d", msg.ArticleID, msg.UserID)
		}
	}
}
