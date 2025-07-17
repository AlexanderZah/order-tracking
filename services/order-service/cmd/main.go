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
	etaClient "github.com/AlexanderZah/order-tracking/services/order-service/internal/client/etaservice"
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
	pool, err := pgxpool.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("can't create pg pool: %s", err.Error())
	}
	repo := orderRepo.New(pool)

	etacl, err := etaClient.New("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to ETA service: %v", err)
	}

	order_ucase := orderUcase.New(repo, etacl)
	defer producer.Close()
	// инициализируем роутер
	router, err := handler.Router(context.Background(), logger, cfg, producer, order_ucase)
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
