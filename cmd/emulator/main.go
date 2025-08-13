package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	k "github.com/unkabas/wb-L0/internal/kafka"
	"github.com/unkabas/wb-L0/models"
	"log"
	"math/rand"
	"time"
)

const (
	topic         = "order"
	numMessages   = 10  // Количество сообщений
	invalidChance = 0.2 // Шанс отправки невалидного JSON
)

var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	p, err := k.NewProducer(address)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numMessages; i++ {
		var data []byte
		if rand.Float64() < invalidChance {
			data = []byte(generateInvalidJSON())
		} else {
			order := generateValidOrder()
			data, err = json.Marshal(order)
			if err != nil {
				log.Printf("Ошибка сериализации заказа: %v", err)
				continue
			}
		}

		key := fmt.Sprintf("key-%d", i)
		if err := p.Produce(string(data), topic, key); err != nil {
			log.Printf("Failed to produce message: %v", err)
		} else {
			log.Printf("Successfully produced message with key: %s\n", key)
		}

		time.Sleep(time.Second)
	}
}

// генерирует разнообразный валидный заказ
func generateValidOrder() models.Order {
	faker := gofakeit.New(0)
	orderUID := faker.UUID()
	numItems := faker.Number(1, 3) // 1-3 товара
	items := make([]models.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = models.Item{
			ChrtID:      faker.Number(1000, 9999),
			TrackNumber: faker.UUID(),
			Price:       faker.Number(100, 1000),
			Rid:         faker.UUID(),
			Name:        faker.ProductName(),
			Sale:        faker.Number(0, 50),
			Size:        faker.Word(),
			TotalPrice:  faker.Number(100, 1000),
			NmID:        faker.Number(1000, 9999),
			Brand:       faker.Company(),
			Status:      faker.Number(100, 200),
		}
	}
	return models.Order{
		OrderUID:    orderUID,
		TrackNumber: faker.UUID(),
		Entry:       faker.Word(),
		Delivery: models.Delivery{
			Name:    faker.Name(),
			Phone:   faker.PhoneFormatted(),
			Zip:     faker.Zip(),
			City:    faker.City(),
			Address: faker.Street(),
			Region:  faker.State(),
			Email:   faker.Email(),
		},
		Payment: models.Payment{
			Transaction:  faker.UUID(),
			RequestID:    faker.UUID(),
			Currency:     faker.CurrencyShort(),
			Provider:     "wbpay",
			Amount:       faker.Number(100, 5000),
			PaymentDt:    faker.DateRange(time.Now().Add(-24*time.Hour), time.Now()).Unix(),
			Bank:         faker.Company(),
			DeliveryCost: faker.Number(100, 1000),
			GoodsTotal:   faker.Number(100, 5000),
			CustomFee:    faker.Number(0, 500),
		},
		Items:             items,
		Locale:            faker.LanguageAbbreviation(),
		InternalSignature: "",
		CustomerID:        faker.UUID(),
		DeliveryService:   "meest",
		Shardkey:          fmt.Sprintf("%d", faker.Number(1, 10)),
		SmID:              faker.Number(1, 100),
		DateCreated:       faker.DateRange(time.Now().Add(-24*time.Hour), time.Now()),
		OofShard:          fmt.Sprintf("%d", faker.Number(1, 10)),
	}
}

// generateInvalidJSON генерирует невалидный JSON
func generateInvalidJSON() string {
	invalidTypes := []string{
		"просто текст, не JSON",
		"{\"order_uid\": \"123\",",
		"{order_uid: 123}",
		"{\"order_uid\": \"\", \"track_number\": \"123\"}",
		"{\"order_uid\": \"123\", \"delivery\": {}}",
		"{\"order_uid\": \"123\", \"items\": []}",
	}
	return invalidTypes[rand.Intn(len(invalidTypes))]
}
