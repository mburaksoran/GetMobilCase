package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/repository/engines"
)

type OrderRepository struct {
	mySqlClient *sql.DB
}

func NewOrderRepository() *OrderRepository {
	sqlDbEngine := engines.GetSqlDbEngine()
	return &OrderRepository{
		mySqlClient: sqlDbEngine.Client,
	}
}

func (r *OrderRepository) CreateOrder(order models.Order) error {
	queryString := `INSERT INTO orders(id,user_id,product_id,ordered_count,price) VALUES (?,?,?,?,?)`

	args := prepareCreateArgsForOrders(order)

	stmt, err := r.mySqlClient.Prepare(queryString)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to prepare query for insert order:%s", err.Error()))
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
		return nil, errors.New(fmt.Sprintf("failed to prepare query for get order:%s", err.Error()))
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
				break
			}
			fmt.Println(err)

		}
	}

	return &sqlResult, nil
}

func (r *OrderRepository) DeleteOrder(id int) error {
	_, err := r.mySqlClient.Query("SELECT * FROM orders WHERE id = ?", id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to prepare query for delete order:%s", err.Error()))
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
