package kafka

import (
	"encoding/json"
	"gingorm/models"
	"gingorm/service/dto"
	"log"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type LikeConsumer struct {
	db    *gorm.DB
	topic string
}

func NewLikeConsumer(db *gorm.DB, topic string) *LikeConsumer {
	return &LikeConsumer{db: db, topic: topic}
}

func (c *LikeConsumer) Start(brokers []string) error {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()

	partitions, _ := consumer.Partitions(c.topic)
	for _, partition := range partitions {
		pc, _ := consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				var m dto.LikeMessage
				if err := json.Unmarshal(msg.Value, &m); err != nil {
					log.Println("[ERROR] invalid like message:", err)
					continue
				}

				switch m.Action {
				case "like":
					// 幂等插入
					if err := c.db.Exec(`
						INSERT INTO likes (user_id, target_id, target_type, created_at, updated_at)
						VALUES (?, ?, ?, NOW(), NOW())
						ON DUPLICATE KEY UPDATE updated_at = NOW()
					`, m.UserID, m.TargetID, m.TargetType).Error; err != nil {
						log.Printf("[Consumer] failed to insert like: %v", err)
					} else {
						log.Printf("[Consumer] like inserted: %+v", m)
					}

				case "unlike":
					// 删除操作
					if err := c.db.Where("user_id = ? AND target_id = ? AND target_type = ?", m.UserID, m.TargetID, m.TargetType).
						Delete(&models.Like{}).Error; err != nil {
						log.Printf("[Consumer] failed to delete like: %v", err)
					} else {
						log.Printf("[Consumer] like deleted: %+v", m)
					}
				}
			}
		}(pc)
	}
	select {} // 阻塞主线程
}
