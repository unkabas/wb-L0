package cache

import (
	"github.com/unkabas/wb-L0/models"
	"gorm.io/gorm"
	"log"
	"sync"
)

const cashSize = 10

type Cache struct {
	data    map[string]models.Order
	order   []string
	maxSize int
	mu      sync.RWMutex
}

// создаёт новый кэш с ограничением размера
func NewCache() *Cache {
	return &Cache{
		data:    make(map[string]models.Order),
		order:   make([]string, 0, cashSize),
		maxSize: cashSize,
	}
}

// загружает последние maxSize заказов из БД в кэш
func (c *Cache) Init(db *gorm.DB) error {
	var orders []models.Order
	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").
		Order("date_created DESC").Limit(c.maxSize).Find(&orders).Error; err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.data[order.OrderUID] = order
		c.order = append(c.order, order.OrderUID)
		log.Printf("Loaded order to cache: %s\n", order.OrderUID)
	}
	log.Printf("Cache initialized with %d orders from DB\n", len(orders))
	return nil
}

// добавляет и удаляет заказы
func (c *Cache) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Добавляем новый заказ
	c.data[order.OrderUID] = order
	c.order = append(c.order, order.OrderUID)
	log.Printf("Added order to cache: %s\n", order.OrderUID)
	// Если кэш полон, удаляем самый старый заказ
	if len(c.data) >= c.maxSize {
		oldOrderUID := c.order[0]
		c.order = c.order[1:]
		delete(c.data, oldOrderUID)
		log.Printf("Removed oldest order from cache: %s\n", oldOrderUID)
	}
}

// возвращает заказ по order_uid, возвращает false, если не найден
func (c *Cache) Get(orderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.data[orderUID]
	return order, exists
}
