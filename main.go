package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"userservice/events"
	"userservice/kafka"
	"userservice/redis"
)

func main() {

	redis.NewRedisClient()

	// Kafka Producer
	producer, err := kafka.NewKafkaProducer()
	if err != nil {
		log.Fatalf("failed to init kafka producer: %v", err)
	}

	// Kafka Consumer (async)
	go kafka.StartKafkaConsumer()

	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Produce order event
	r.Post("/produce", func(w http.ResponseWriter, r *http.Request) {

		var req events.OrderEvent
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", 400)
			return
		}

		data, _ := json.Marshal(req)

		if err := producer.Produce("orders", string(data)); err != nil {
			log.Println(err)
			http.Error(w, "failed to send message", 500)
			return
		}

		w.Write([]byte("sent"))
	})

	// Get trending products
	r.Get("/trending", func(w http.ResponseWriter, r *http.Request) {

		limitStr := r.URL.Query().Get("limit")
		limit, _ := strconv.Atoi(limitStr)
		if limit == 0 {
			limit = 10
		}

		items, err := redis.GetTopTrending(limit)
		if err != nil {
			http.Error(w, "redis error", 500)
			return
		}

		json.NewEncoder(w).Encode(items)
	})

	log.Println("server running on :8080")
	http.ListenAndServe("localhost:8080", r)
}

func UtilsResponseJson(w http.ResponseWriter, status int, payload map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
