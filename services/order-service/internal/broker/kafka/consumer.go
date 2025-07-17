package kafka

import (
	"context"
	"encoding/json"
	"log"

	order_ucase "github.com/AlexanderZah/order-tracking/services/order-service/internal/app/usecase/order"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/event"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	reader *kafka.Reader
	uCase  *order_ucase.Usecase
	log    logrus.FieldLogger
}

func NewConsumer(brokers []string, topic string, uCase *order_ucase.Usecase, log logrus.FieldLogger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: "eta-service",
			Topic:   topic,
		}),
		uCase: uCase,
		log:   log}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading msg: %v", err)
			return err
		}

		var orderUp event.OrderETAUpdated
		if err := json.Unmarshal(m.Value, &orderUp); err != nil {
			log.Printf("error unmarshaling: %v", err)
			return err
		}

		c.uCase.Update(ctx, c.log, orderUp.OrderID, &orderUp.ETAMinutes)

	}
}
