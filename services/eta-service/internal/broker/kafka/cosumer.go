package kafka

import (
	"context"
	"encoding/json"
	"log"

	etaClient "github.com/AlexanderZah/order-tracking/services/eta-service/internal/client"
	"github.com/AlexanderZah/order-tracking/services/eta-service/internal/event"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: "eta-service",
			Topic:   topic,
		})}
}

func (c *Consumer) Consume(ctx context.Context, eta *etaClient.Client, producer *Producer) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading msg: %v", err)
			return err
		}

		var orderCr event.OrderCreated
		if err := json.Unmarshal(m.Value, &orderCr); err != nil {
			log.Printf("error unmarshaling: %v", err)
			return err
		}

		eta, err := eta.GetETA(ctx, orderCr.Address)
		if err != nil {
			log.Printf("error getting ETA: %v", err)
			return err
		}

		log.Printf("ETA for order %s is %d minutes", orderCr.OrderID, eta)
		err = producer.PublishETAUpdated(ctx, "eta-updated", event.OrderETAUpdated{
			OrderID:    orderCr.OrderID,
			ETAMinutes: int(eta),
		})
		if err != nil {
			log.Printf("error publishing ETA updated: %v", err)
		}
	}
}
