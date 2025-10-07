package main

import (
	"gingorm/config"
	"gingorm/global"
	"gingorm/kafka"
	"gingorm/router"
	"log"
)

func main() {
	config.InitConfig()
	//kafka.InitKafkaWriter([]string{"localhost:9092"}, "like-topic")
	//defer kafka.CloseKafkaWriter()
	brokers := []string{"localhost:9092"} // Kafka broker 地址
	topic := "like_events"
	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	go func() {
		c := kafka.NewLikeConsumer(global.DB, topic)
		if err := c.Start(brokers); err != nil {
			log.Fatalf("Kafka consumer failed: %v", err)
		}
	}()
	r := router.SetupRouter(global.DB, global.RedisDB, producer)
	//tasks.StartLikesSyncTicker()
	//go kafka.StartLikeConsumer([]string{"localhost:9092"}, "like-topic", "like-group")
	r.Run(":8080")
}
