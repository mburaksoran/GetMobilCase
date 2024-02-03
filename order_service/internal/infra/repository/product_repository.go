package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/repository/engines"
)

type ProductRepository struct {
	mySqlClient *sql.DB
}

func NewProductRepository() *ProductRepository {
	sqlDbEngine := engines.GetSqlDbEngine()
	return &ProductRepository{
		mySqlClient: sqlDbEngine.Client,
	}
}

func (r *ProductRepository) CreateProduct(product models.Product) error {
	queryString := `INSERT INTO products(id,name,stock_count,price) VALUES (?,?,?,?)`

	args := prepareCreateArgsForProduct(product)

	stmt, err := r.mySqlClient.Prepare(queryString)
	if err != nil {
		return errors.New("failed to prepare query for insert product")
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		args...,
	)
	return err
}

func (r *ProductRepository) GetProduct(id int) (*models.Product, error) {
	row, err := r.mySqlClient.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, errors.New("failed to prepare query for insert product")
	}
	var sqlResult models.Product

	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(
			&sqlResult.ID,
			&sqlResult.Name,
			&sqlResult.StockCount,
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

func (r *ProductRepository) DeleteProduct(id int) error {
	_, err := r.mySqlClient.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return errors.New("failed to query for delete product")
	}
	return nil
}

func (r *ProductRepository) UpdateStockCount(value int, productId int) error {
	tx, err := r.mySqlClient.Begin()
	if err != nil {
		fmt.Println(err)
	}

	_, err = tx.Exec("UPDATE products SET stock_count = stock_count - ? WHERE id = ?", value, productId)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		return errors.New("failed to query for update stock count product")
	}
	return nil
}

func prepareCreateArgsForProduct(product models.Product) []interface{} {
	queryArgs := make([]interface{}, 0)
	queryArgs = append(queryArgs, product.ID)
	queryArgs = append(queryArgs, product.Name)
	queryArgs = append(queryArgs, product.StockCount)
	queryArgs = append(queryArgs, product.Price)
	return queryArgs
}
