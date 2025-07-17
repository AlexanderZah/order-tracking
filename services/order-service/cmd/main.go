package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	config "github.com/AlexanderZah/order-tracking/services/order-service/config"
	orderUcase "github.com/AlexanderZah/order-tracking/services/order-service/internal/app/usecase/order"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/broker/kafka"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/handler"
	orderRepo "github.com/AlexanderZah/order-tracking/services/order-service/internal/repository/order"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logrus.New()
	brokers := strings.Split(cfg.KafkaAddr, ",")
	producer := kafka.New(brokers, "order.created")
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("can't create pg pool: %s", err.Error())
	}
	repo := orderRepo.New(pool)

	order_ucase := orderUcase.New(repo)
	defer producer.Close()
	// инициализируем роутер
	router, err := handler.Router(ctx, logger, cfg, producer, order_ucase)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("Server started on port %s", cfg.HTTPPort)

	// запускаем HTTP-сервер
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
