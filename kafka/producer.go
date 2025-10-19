package kafka

import (
	"context"
	"encoding/json"
	"gingorm/service/dto"
	kafkago "github.com/segmentio/kafka-go"
	"log"
	"time"
)

type LikeProducer struct {
	writer *kafkago.Writer
}

func NewLikeProducer(brokers []string, topic string) *LikeProducer {
	return &LikeProducer{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafkago.LeastBytes{},
			BatchSize:    1,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *LikeProducer) SendLike(msg dto.LikeMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[Producer] JSON marshal error: %v", err)
		return
	}

	err = p.writer.WriteMessages(context.Background(),
		kafkago.Message{
			Key:   []byte(msg.Action),
			Value: data,
		},
	)
	if err != nil {
		log.Printf("[Producer] send message failed: %v", err)
	}
}

func (p *LikeProducer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("[Producer] close error: %v", err)
	}
}
