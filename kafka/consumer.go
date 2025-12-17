package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"userservice/events"
	"userservice/redis"

	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "order-group-1",
		
	})

	defer r.Close()
	log.Println("Kafka consumer started...")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading msg: %v", err)
			continue
		}

		fmt.Printf("received => %s\n", string(msg.Value))

		var event events.OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err == nil {
			// increment trending product count
			redis.IncrementTrending(fmt.Sprintf("%d", event.ProductId))
		}
	}
}
