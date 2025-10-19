package kafka

import (
	"context"
	"encoding/json"
	"gingorm/models"
	"gingorm/service/dto"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"log"
)

type LikeConsumer struct {
	db     *gorm.DB
	reader *kafka.Reader
}

func NewLikeConsumer(db *gorm.DB, brokers []string, topic, groupID string) *LikeConsumer {
	return &LikeConsumer{
		db: db,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 1,
			MaxBytes: 1e6, // 1MB
		}),
	}
}

func (c *LikeConsumer) Start() {
	go func() {
		for {
			m, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[Consumer] read error: %v", err)
				continue
			}

			var msg dto.LikeMessage
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("[Consumer] invalid message: %v", err)
				continue
			}

			switch msg.Action {
			case "like":
				err = c.db.Exec(`
					INSERT INTO likes (user_id, target_id, target_type, created_at, updated_at)
					VALUES (?, ?, ?, NOW(), NOW())
					ON DUPLICATE KEY UPDATE updated_at = NOW()
				`, msg.UserID, msg.TargetID, msg.TargetType).Error
			case "unlike":
				err = c.db.Where("user_id=? AND target_id=? AND target_type=?", msg.UserID, msg.TargetID, msg.TargetType).
					Delete(&models.Like{}).Error
			}
			if err != nil {
				log.Printf("[Consumer] db operation failed: %v", err)
			} else {
				log.Printf("[Consumer] handled action=%s user=%d target=%d", msg.Action, msg.UserID, msg.TargetID)
			}
		}
	}()
}

func (c *LikeConsumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Printf("[Consumer] close error: %v", err)
	}
}
