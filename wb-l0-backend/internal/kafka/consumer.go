package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/unkabas/wb-L0/internal/cache"
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
	cache    *cache.Cache
}

func NewConsumer(address []string, topic, groupID string, db *gorm.DB, cache *cache.Cache) (*Consumer, error) {
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
		cache:    cache,
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

			// Проверка пустого сообщения
			if len(msg.Value) == 0 {
				log.Println("We got empty message, continue")
				if _, err := c.consumer.CommitMessage(msg); err != nil {
					log.Printf("Error commited empty message: %v", err)
				}
				continue
			}

			if err := c.processMessage(msg.Value); err != nil {
				log.Printf("Failed to process message %v", err)
				continue
			}

			if _, err := c.consumer.CommitMessage(msg); err != nil {
				log.Printf("Failed to commit offset: %v", err)
			} else {
				log.Printf("Message committed successfully for message\n")
			}
		}
	}
}

func (c *Consumer) processMessage(data []byte) (err error) {
	// Защита от паник
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in process Message: %v, message: %q", r, string(data))
			err = fmt.Errorf("panic in processMessage: %v", r)
		}
	}()

	// Парсинг JSON
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("unmarshal error for message %q: %w", string(data), err)
	}

	// Валидация ключевых полей
	if order.OrderUID == "" {
		return fmt.Errorf("invalid order: order_uid is empty, message: %q", string(data))
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		return fmt.Errorf("invalid order: delivery name or phone is empty, order_uid: %s", order.OrderUID)
	}
	if order.Payment.Transaction == "" {
		return fmt.Errorf("invalid order: payment transaction is empty, order_uid: %s", order.OrderUID)
	}

	// Retry-логика для сохранения в БД
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		tx := c.db.Begin()
		// Откат при панике внутри транзакции
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Printf("Panic during transaction for order %s, attempt %d: %v", order.OrderUID, attempt, r)
			}
		}()

		// Сохранение заказа
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			if attempt == maxRetries {
				return fmt.Errorf("create order failed after %d attempts, order_uid: %s: %w", maxRetries, order.OrderUID, err)
			}
			log.Printf("Attempt %d failed for order %s: %v, retrying...", attempt, order.OrderUID, err)
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		// Коммит транзакции
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("commit failed for order %s: %w", order.OrderUID, err)
		}

		c.cache.Set(order)
		log.Printf("Successfully saved order %s and added to cache\n", order.OrderUID)
		return nil

	}
	return fmt.Errorf("create order failed after %d attempts, order_uid: %s", maxRetries, order.OrderUID)
}

func (c *Consumer) Stop() {
	close(c.stopChan)
	if err := c.consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}
