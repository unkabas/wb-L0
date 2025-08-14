package api

import (
	"encoding/json"
	"github.com/unkabas/wb-L0/internal/cache"
	"github.com/unkabas/wb-L0/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func OrderHandler(cache *cache.Cache, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method != http.MethodGet {
			http.Error(w, "Only get", 405)
			return
		}

		orderUID := strings.TrimPrefix(r.URL.Path, "/order/")
		if orderUID == "" {
			http.Error(w, "No order_uid", 400)
			return
		}

		// Проверяем кэш
		order, exists := cache.Get(orderUID)
		if !exists {
			// Ищем в БД
			var dbOrder models.Order
			if err := db.Preload("Delivery").Preload("Payment").Preload("Items").
				Where("order_uid = ?", orderUID).First(&dbOrder).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					http.Error(w, "Order not found", 404)
					return
				}
				log.Printf("Error in db with order_uid %s: %v", orderUID, err)
				http.Error(w, "Something went wrong", 500)
				return
			}
			cache.Set(dbOrder)
			order = dbOrder
			log.Printf("Order %s set in cache", orderUID)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Printf("JSON error with order_uid %s: %v", orderUID, err)
			http.Error(w, "Something went wrong", 500)
			return
		}
	}
}
