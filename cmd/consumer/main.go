package main

import (
	"github.com/unkabas/wb-L0/internal/config"
	"github.com/unkabas/wb-L0/internal/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	topic   = "order"
	address = "localhost:9091,localhost:9092,localhost:9093"
)

func init() {
	config.LoadEnvs()
	config.ConnectDB()
}

func main() {
	consumer, err := kafka.NewConsumer(
		[]string{"localhost:9091", "localhost:9092", "localhost:9093"},
		topic,
		"orders-group",
		config.DB,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Stop()

	go func() {
		log.Println("Starting Kafka consumer...")
		consumer.Start()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer...")
}
