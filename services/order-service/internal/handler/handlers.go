package handler

import (
	"context"
	"fmt"

	"github.com/AlexanderZah/order-tracking/services/order-service/config"
	order_ucase "github.com/AlexanderZah/order-tracking/services/order-service/internal/app/usecase/order"
	create_order_handler "github.com/AlexanderZah/order-tracking/services/order-service/internal/handler/order/create"
	get_orders_handler "github.com/AlexanderZah/order-tracking/services/order-service/internal/handler/order/get"
	orderRepo "github.com/AlexanderZah/order-tracking/services/order-service/internal/repository/order"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	orderRoute  = "/order"
	ordersRoute = "/orders"
)

// Router register necessary routes and returns an instance of a router.
func Router(ctx context.Context, log logrus.FieldLogger, config *config.Config) (*mux.Router, error) {
	r := mux.NewRouter()

	// echo

	pool, err := pgxpool.Connect(context.Background(), config.DBURL)
	if err != nil {
		return nil, fmt.Errorf("can't create pg pool: %s", err.Error())
	}
	repo := orderRepo.New(pool)
	order_ucase := order_ucase.New(repo)

	createOrderHandleFunc := create_order_handler.New(order_ucase, log).Create(ctx).ServeHTTP
	// create order
	r.HandleFunc(orderRoute, createOrderHandleFunc).Methods("POST")
	getOrderHandlerFunc := get_orders_handler.New(order_ucase, log).Get(ctx).ServeHTTP
	// get orders
	r.HandleFunc(ordersRoute, getOrderHandlerFunc).Methods("GET")

	return r, nil
}
