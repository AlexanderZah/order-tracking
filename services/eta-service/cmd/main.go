package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AlexanderZah/order-tracking/services/eta-service/config"
	pb "github.com/AlexanderZah/order-tracking/services/eta-service/gen/go/etaservice/v1"
	kafkaConsumer "github.com/AlexanderZah/order-tracking/services/eta-service/internal/broker/kafka/consumer"
	etaClient "github.com/AlexanderZah/order-tracking/services/eta-service/internal/client"
	etaServer "github.com/AlexanderZah/order-tracking/services/eta-service/internal/server"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Основной контекст с graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	pb.RegisterETAServiceServer(grpcServer, etaServer.New())

	// Инициализация клиента ETA
	etaCl, err := etaClient.New(cfg.ETAService)
	if err != nil {
		log.Fatalf("can't create ETA client: %v", err)
	}

	// Инициализация Kafka consumer
	fmt.Println(cfg.KafkaAddr)
	consumer := kafkaConsumer.New(strings.Split(cfg.KafkaAddr, ","), "order.created")

	// Запуск consumer в отдельной горутине с обработкой ошибок
	consumerDone := make(chan error, 1)
	go func() {
		defer close(consumerDone)
		if err := consumer.Consume(ctx, etaCl); err != nil {
			consumerDone <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	// Запуск gRPC сервера с graceful shutdown
	listener, err := net.Listen("tcp", cfg.ETAService)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server starting on %s", cfg.ETAService)

	// Канал для ошибок сервера
	serverDone := make(chan error, 1)
	go func() {
		defer close(serverDone)
		if err := grpcServer.Serve(listener); err != nil {
			serverDone <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Ожидание сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		cancel() // Отправляем сигнал отмены в контекст
		grpcServer.GracefulStop()
	case err := <-consumerDone:
		log.Printf("Consumer stopped: %v", err)
		grpcServer.GracefulStop()
	case err := <-serverDone:
		log.Printf("gRPC server stopped: %v", err)
		cancel()
	}

	log.Println("Shutdown completed")
}
