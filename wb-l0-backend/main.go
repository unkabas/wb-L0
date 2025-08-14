package main

import (
	"flag"
	"github.com/unkabas/wb-L0/cmd/migration"
	"github.com/unkabas/wb-L0/internal/api"
	"github.com/unkabas/wb-L0/internal/cache"
	"github.com/unkabas/wb-L0/internal/config"
	"github.com/unkabas/wb-L0/internal/kafka"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var migrate = flag.Bool("m", false, "migration")

const (
	topic = "order"
)

var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func init() {
	config.LoadEnvs()
	config.ConnectDB()
}
func main() {
	flag.Parse()
	if *migrate {
		migration.Migration()
	}

	//инициализация кэша
	c := cache.NewCache()
	if err := c.Init(config.DB); err != nil {
		log.Printf("Failed to initialize cache: %v", err)
	}

	groupID := "orders-group"
	//запуск consumer
	consumer, err := kafka.NewConsumer(
		address,
		topic,
		groupID,
		config.DB,
		c,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	go func() {
		log.Println("Starting Kafka consumer...")
		consumer.Start()
	}()
	defer consumer.Stop()

	http.HandleFunc("/order/", api.OrderHandler(c, config.DB))
	go func() {
		log.Println("Listening port :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	//ждём сигнал завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}
