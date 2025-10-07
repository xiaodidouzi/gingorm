package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"log"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.WaitForLocal
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{asyncProducer: producer, topic: topic}, nil
}

func (p *Producer) SendMessage(ctx context.Context, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	}:
		log.Printf("[INFO] Kafka message sent: %+v", msg)
		return nil
	case <-ctx.Done():
		return errors.New("send kafka message canceled or timeout: " + ctx.Err().Error())
	}
}
