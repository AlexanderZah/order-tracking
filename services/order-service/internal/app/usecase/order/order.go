package order

import (
	"context"
	"fmt"

	etaClient "github.com/AlexanderZah/order-tracking/services/order-service/internal/client/etaservice"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/repository/entity/order"
	orderRepo "github.com/AlexanderZah/order-tracking/services/order-service/internal/repository/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Usecase struct {
	repo orderRepo.OrderRepo
	eta  *etaClient.Client
}

// New gives Usecase.
func New(orderRepo orderRepo.OrderRepo, eta *etaClient.Client) *Usecase {
	return &Usecase{repo: orderRepo,
		eta: eta}
}

func (uc *Usecase) Save(ctx context.Context, log logrus.FieldLogger, order *order.Order) error {
	if err := uc.repo.Save(ctx, log, order); err != nil {

		return err
	}
	etaMinutes, err := uc.eta.GetETA(ctx, order.DeliveryAddress)
	if err != nil {
		return err
	}
	eta := int(etaMinutes)
	order.ETAMinutes = &eta

	return nil
}

func (uc *Usecase) Get(ctx context.Context, log logrus.FieldLogger, IDs []uuid.UUID) ([]order.Order, error) {
	ordersMap, err := uc.repo.Get(ctx, log, IDs)
	if err != nil {

		return nil, fmt.Errorf("err from orders_repository: %s", err.Error())
	}

	result := make([]order.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		result = append(result, order)
	}

	return result, nil
}
