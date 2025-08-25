package main

import (
	"awesomeProject/config"
	"awesomeProject/kafka"
	"awesomeProject/router"
	"awesomeProject/tasks"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.InitConfig()
	kafka.InitKafkaWriter([]string{"localhost:9092"}, "like-topic")
	defer kafka.CloseKafkaWriter()
	r := router.SetupRouter()
	tasks.StartLikesSyncTicker()
	go kafka.StartLikeConsumer([]string{"localhost:9092"}, "like-topic", "like-group")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
