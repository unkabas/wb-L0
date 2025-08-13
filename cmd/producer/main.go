package main

import (
	"encoding/json"
	"github.com/google/uuid"
	k "github.com/unkabas/wb-L0/internal/kafka"
	"log"
)

const (
	topic = "order"
)

func main() {
	p, err := k.NewProducer([]string{"localhost:9091", "localhost:9092", "localhost:9093"})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	for i := 0; i < 3; i++ {
		orderUID := uuid.New().String()
		order := map[string]interface{}{
			"order_uid":    orderUID,
			"track_number": "WBILMTESTTRACK",
			"entry":        "WBIL",
			"delivery": map[string]interface{}{
				"name":    "Test Testov",
				"phone":   "+9720000000",
				"zip":     "2639809",
				"city":    "Kiryat Mozkin",
				"address": "Ploshad Mira 15",
				"region":  "Kraiot",
				"email":   "test@gmail.com",
			},
			"payment": map[string]interface{}{
				"transaction":   orderUID,
				"request_id":    "",
				"currency":      "USD",
				"provider":      "wbpay",
				"amount":        1817,
				"payment_dt":    1637907727,
				"bank":          "alpha",
				"delivery_cost": 1500,
				"goods_total":   317,
				"custom_fee":    0,
			},
			"items": []map[string]interface{}{
				{
					"chrt_id":      9934930,
					"track_number": "WBILMTESTTRACK",
					"price":        453,
					"rid":          uuid.New().String(),
					"name":         "Mascaras",
					"sale":         30,
					"size":         "0",
					"total_price":  317,
					"nm_id":        2389212,
					"brand":        "Vivienne Sabo",
					"status":       202,
				},
			},
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        "test",
			"delivery_service":   "meest",
			"shardkey":           "9",
			"sm_id":              99,
			"date_created":       "2021-11-26T06:22:19Z",
			"oof_shard":          "1",
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order to JSON: %v", err)
			continue
		}

		key := uuid.New().String()
		if err := p.Produce(string(jsonData), topic, key); err != nil {
			log.Printf("Failed to produce message: %v", err)
		} else {
			log.Printf("Successfully produced message with key: %s, order_uid: %s\n", key, orderUID)
		}
	}
}
