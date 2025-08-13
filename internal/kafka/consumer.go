package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/unkabas/wb-L0/models"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type Consumer struct {
	consumer *kafka.Consumer
	db       *gorm.DB
	stopChan chan struct{}
}

func NewConsumer(address []string, topic, groupID string, db *gorm.DB) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer: %v", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		c.Close()
		return nil, fmt.Errorf("error subscribing to topic %s: %v", topic, err)
	}

	return &Consumer{
		consumer: c,
		db:       db,
		stopChan: make(chan struct{}),
	}, nil
}

func (c *Consumer) Start() {
	log.Println("Starting Kafka consumer")
	defer log.Println("Kafka consumer stopped")

	for {
		select {
		case <-c.stopChan:
			return
		default:
			msg, err := c.consumer.ReadMessage(5 * time.Second)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					log.Println("Polling timeout, no new messages")
					continue
				}
				log.Printf("Consumer error: %v", err)
				continue
			}

			if err := c.processMessage(msg.Value); err != nil {
				log.Printf("Failed to process message: %v", err)
				continue
			}

			if _, err := c.consumer.CommitMessage(msg); err != nil {
				log.Printf("Failed to commit offset: %v", err)
			} else {
				log.Println("Message committed successfully")
			}
		}
	}
}

func (c *Consumer) processMessage(data []byte) error {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic during transaction: %v", r)
		}
	}()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create order failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	log.Printf("Successfully saved order %s\n", order.OrderUID)
	return nil
}

func (c *Consumer) Stop() {
	close(c.stopChan)
	if err := c.consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}
