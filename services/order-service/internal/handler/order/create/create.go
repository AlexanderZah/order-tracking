package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	order_ucase "github.com/AlexanderZah/order-tracking/services/order-service/internal/app/usecase/order"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/broker/kafka"
	etaClient "github.com/AlexanderZah/order-tracking/services/order-service/internal/client/etaservice"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/event"
	"github.com/AlexanderZah/order-tracking/services/order-service/internal/repository/entity/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Requst validation errors.
var ErrInvalidUserID = errors.New("invalid user ID")
var ErrInvalidAmount = errors.New("invalid price")
var ErrInvalidPaymentType = errors.New("invalid payment type")
var ErrEmptyItems = errors.New("items can't be empty")
var ErrInvalidItemID = errors.New("invalid service id")

// Handler creates orders
type Handler struct {
	uCase    *order_ucase.Usecase
	log      logrus.FieldLogger
	producer *kafka.Producer
}

// New gives Handler.
func New(
	uCase *order_ucase.Usecase,
	log logrus.FieldLogger,
	producer *kafka.Producer,
) *Handler {
	return &Handler{
		uCase:    uCase,
		log:      log,
		producer: producer,
	}
}

type OrderIn struct {
	UserID          uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	ETAMinutes      *int      `json:"eta_minutes,omitempty" db:"eta_minutes"`
	DeliveryAddress string    `json:"delivery_address" db:"delivery_address" validate:"required,min=10,max=500"`
}

func (in OrderIn) OrderFromDTO() order.Order {
	return order.Order{
		UserID:          in.UserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ETAMinutes:      in.ETAMinutes,
		DeliveryAddress: in.DeliveryAddress,
		Status:          order.StatusCreated,
	}
}

func (h Handler) validateReq(in *OrderIn) error {
	_ = in
	return nil
}

func (h Handler) Create(ctx context.Context) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// prepare dto to parse request
		in := &OrderIn{}
		// parse req body to dto
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			h.log.Errorf("can't parse req: %s", err.Error())
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}

		// check that request valid
		err = h.validateReq(in)
		if err != nil {
			h.log.Errorf("bad req: %v: %s", in, err.Error())
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		order := in.OrderFromDTO()
		err = h.uCase.Save(ctx, h.log, &order)
		if err != nil {
			h.log.Errorf("can't create order: %v: %s", order, err.Error())
			http.Error(w, "can't create order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		createdEvent := event.OrderCreated{
			OrderID: order.ID,
			UserID:  order.UserID,
			Address: order.DeliveryAddress,
		}

		err = h.producer.Publish(ctx, order.ID.String(), createdEvent)
		if err != nil {
			h.log.Errorf("can't publish order.created: %s", err)
		}

		eta, err := etaClient.GetETA(ctx, order.DeliveryAddress)
		order.ETAMinutes = &eta

		w.Header().Set("Content-Type", "application/json")

		m := make(map[string]interface{})
		m["success"] = "ok"
		json.NewEncoder(w).Encode(m)

	}
	return http.HandlerFunc(fn)
}
