package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type LikeMessage struct {
	ArticleID int `json:"article_id"`
	UserID    int `json:"user_id"`
}

var Writer *kafka.Writer

func InitKafkaWriter(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func SendLikeMessage(articleID, userID int) error {
	msg := LikeMessage{
		ArticleID: articleID,
		UserID:    userID,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
		log.Printf("Kafka send failed: ArticleID=%d UserID=%d, err=%v", articleID, userID, err)
		return err
	}

	log.Printf("Kafka message sent: ArticleID=%d UserID=%d", articleID, userID)
	return nil
}

func CloseKafkaWriter() {
	if Writer != nil {
		_ = Writer.Close()
	}
}
