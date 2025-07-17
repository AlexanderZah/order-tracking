package order

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	order_ucase "github.com/AlexanderZah/order-tracking/services/order-service/internal/app/usecase/order"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var ErrEmptyOrderId = errors.New("no order id passed")

type Handler struct {
	uCase *order_ucase.Usecase
	log   logrus.FieldLogger
}

func New(
	uCase *order_ucase.Usecase,
	log logrus.FieldLogger,
) *Handler {
	return &Handler{
		uCase: uCase,
		log:   log,
	}
}

type UpdateOrderIn struct {
	ID         uuid.UUID `json:"id"`
	EtaMinutes *int      `json:"eta_minutes"`
}

func (h Handler) validateReq(in *UpdateOrderIn) error {
	if len(in.ID) == 0 {
		return ErrEmptyOrderId
	}

	return nil
}

func (h Handler) Update(ctx context.Context) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// prepare dto to parse request
		in := &UpdateOrderIn{}
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

		err = h.uCase.Update(ctx, h.log, in.ID, in.EtaMinutes)
		if err != nil {
			h.log.Errorf("can't update order: %s", err.Error())
			http.Error(w, "can't update order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
	return http.HandlerFunc(fn)
}
