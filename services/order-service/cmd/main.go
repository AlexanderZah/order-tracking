package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	config "github.com/AlexanderZah/order-tracking/services/order-service/config"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logrus.New()
	// инициализируем роутер
	router, err := handler.Router(context.Background(), logger, cfg)
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
