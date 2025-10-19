package main

import (
	"gingorm/config"
	"gingorm/global"
	"gingorm/kafka"
	"gingorm/router"
)

func main() {
	config.InitConfig()

	global.DB = config.InitDB()
	global.RedisDB = config.InitRedis()

	brokers := []string{"localhost:9092"}
	topic := "like_events"
	groupID := "like_group"

	// Kafka Producer
	producer := kafka.NewLikeProducer(brokers, topic)
	defer producer.Close()

	// Kafka Consumer
	consumer := kafka.NewLikeConsumer(global.DB, brokers, topic, groupID)
	go consumer.Start()
	defer consumer.Close()

	// 启动 HTTP 服务
	r := router.SetupRouter(global.DB, global.RedisDB, producer)
	r.Run(":8080")
}
