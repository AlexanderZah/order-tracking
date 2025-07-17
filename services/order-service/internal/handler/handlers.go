package handler

import (
	"context"

	"github.com/AlexanderZah/order-tracking/services/order-service/config"
	order_ucase "github.com/AlexanderZah/order-tracking/services/order-service/internal/app/usecase/order"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/broker/kafka"
	create_order_handler "github.com/AlexanderZah/order-tracking/services/order-service/internal/handler/order/create"
	get_orders_handler "github.com/AlexanderZah/order-tracking/services/order-service/internal/handler/order/get"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	orderRoute  = "/order"
	ordersRoute = "/orders"
)

// Router register necessary routes and returns an instance of a router.
func Router(ctx context.Context, log logrus.FieldLogger, config *config.Config, producer *kafka.Producer, useCase *order_ucase.Usecase) (*mux.Router, error) {
	r := mux.NewRouter()

	createOrderHandleFunc := create_order_handler.New(useCase, log, producer).Create(ctx).ServeHTTP
	// create order
	r.HandleFunc(orderRoute, createOrderHandleFunc).Methods("POST")
	getOrderHandlerFunc := get_orders_handler.New(useCase, log).Get(ctx).ServeHTTP
	// get orders
	r.HandleFunc(ordersRoute, getOrderHandlerFunc).Methods("GET")

	return r, nil
}
