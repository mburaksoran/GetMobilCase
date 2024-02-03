package repository

import (
	"database/sql"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/repository/engines"
	"go.uber.org/zap"
)

type OrderRepository struct {
	mySqlClient *sql.DB
	logger      *zap.SugaredLogger
}

func NewOrderRepository(lgr *zap.SugaredLogger) *OrderRepository {
	sqlDbEngine := engines.GetSqlDbEngine()
	return &OrderRepository{
		mySqlClient: sqlDbEngine.Client,
		logger:      lgr,
	}
}

func (r *OrderRepository) CreateOrder(order models.Order) error {
	queryString := `INSERT INTO orders(id,user_id,product_id,ordered_count,price) VALUES (?,?,?,?,?)`

	args := prepareCreateArgsForOrders(order)

	stmt, err := r.mySqlClient.Prepare(queryString)
	if err != nil {
		r.logger.Error("failed to prepare query for insert order", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		args...,
	)
	return err
}

func (r *OrderRepository) GetOrder(id int) (*models.Order, error) {
	row, err := r.mySqlClient.Query("SELECT * FROM orders WHERE id = ? LIMIT 1", id)
	if err != nil {
		r.logger.Error("failed to prepare query for get order", err)
		return nil, err
	}
	defer row.Close()
	var sqlResult models.Order

	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(
			&sqlResult.ID,
			&sqlResult.UserID,
			&sqlResult.ProductID,
			&sqlResult.OrderedCount,
			&sqlResult.Price,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				r.logger.Warn(err)
				break
			}
			r.logger.Error(err)

		}
	}

	return &sqlResult, nil
}

func (r *OrderRepository) DeleteOrder(id int) error {
	_, err := r.mySqlClient.Query("SELECT * FROM orders WHERE id = ?", id)
	if err != nil {
		r.logger.Error("failed to prepare query for delete order", err)
		return err
	}
	return nil
}

func prepareCreateArgsForOrders(order models.Order) []interface{} {
	queryArgs := make([]interface{}, 0)
	queryArgs = append(queryArgs, order.ID)
	queryArgs = append(queryArgs, order.UserID)
	queryArgs = append(queryArgs, order.ProductID)
	queryArgs = append(queryArgs, order.OrderedCount)
	queryArgs = append(queryArgs, order.Price)
	return queryArgs
}
