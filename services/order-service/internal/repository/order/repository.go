package order

import (
	"context"
	"fmt"

	"github.com/AlexanderZah/order-tracking/services/order-service/internal/repository/entity/order"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

const (
	// tables
	ordersTable     = "orders"
	itemsTable      = "items"
	orderItemsTable = "order_items"
)

type Repository struct {
	db *pgxpool.Pool
}

type OrderRepo interface {
	Save(ctx context.Context, log logrus.FieldLogger, order *order.Order) error
	Get(ctx context.Context, log logrus.FieldLogger, IDs []uuid.UUID) (map[uuid.UUID]order.Order, error)
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) Save(ctx context.Context, log logrus.FieldLogger, order *order.Order) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("can't create tx: %s", err.Error())
	}

	defer func() {
		// safety rollback if commit didn't happen
		_ = tx.Rollback(ctx)
	}()

	query, args, err := sq.
		Insert(ordersTable).
		Columns("user_id", "status", "created_at", "updated_at", "eta_minutes", "delivery_address").
		Values(order.UserID,
			order.Status,
			order.CreatedAt,
			order.UpdatedAt,
			order.ETAMinutes,
			order.DeliveryAddress).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't build sql:%s", err.Error())
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&order.ID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit tx: %s", err.Error())
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, log logrus.FieldLogger, IDs []uuid.UUID) (map[uuid.UUID]order.Order, error) {
	ordersMap := make(map[uuid.UUID]order.Order, len(IDs))
	or := sq.Or{}
	for _, id := range IDs {
		or = append(or, sq.Eq{"id": id})
	}
	query, args, err := sq.
		Select("id", "user_id", "status", "created_at", "updated_at", "eta_minutes", "delivery_address").
		From(ordersTable).
		Where(or).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("can't build sql:%s", err.Error())
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("can't select orders: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		ord := order.Order{}
		err := rows.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.CreatedAt, &ord.UpdatedAt, &ord.ETAMinutes, &ord.DeliveryAddress)
		if err != nil {
			return nil, fmt.Errorf("can't scan order: %s", err.Error())
		}
		ordersMap[ord.ID] = ord
	}

	return ordersMap, nil
}
