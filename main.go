package main

import (
	"awesomeProject/config"
	"awesomeProject/kafka"
	"awesomeProject/router"
	"awesomeProject/tasks"
)

func main() {
	config.InitConfig()
	kafka.InitKafkaWriter([]string{"localhost:9092"}, "like-topic")
	defer kafka.CloseKafkaWriter()
	r := router.SetupRouter()
	tasks.StartLikesSyncTicker()
	go kafka.StartLikeConsumer([]string{"localhost:9092"}, "like-topic", "like-group")
	r.Run(":8080")
}
