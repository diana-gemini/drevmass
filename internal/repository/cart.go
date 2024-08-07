package repository

import (
	"context"
	"fmt"

	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	cartsTable = "carts"
)

func (r *UserRepository) Add(c context.Context, userId uint, productId int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id) values($1, $2) RETURNING id", cartsTable)
	row := r.db.QueryRow(c, query, userId, productId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) DeleteFromCart(c context.Context, userId uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", cartsTable)
	_, err := r.db.Exec(c, query, userId)
	return err
}

func (r *UserRepository) GetAllFromCart(c context.Context, userId uint) ([]models.GetProductsFromCart, error) {

	query := fmt.Sprintf("SELECT p.id,p.name,p.image,p.price, c.quantity FROM %s p JOIN %s c ON p.id =c.product_id WHERE c.user_id = $1", productsTable, cartsTable)
	rows, err := r.db.Query(c, query, userId)
	var products = []models.GetProductsFromCart{}
	for rows.Next() {
		var product models.GetProductsFromCart
		err := rows.Scan(&product.Id, &product.Name, &product.Image, &product.Price, &product.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, err
}

func (r *UserRepository) GetTotalAmout(c context.Context, userId uint) (float64, error) {
	var count int
	queryCount := fmt.Sprintf("SELECT count(*) FROM %s WHERE user_id = $1", cartsTable)
	rowCount := r.db.QueryRow(c, queryCount, userId)
	if err := rowCount.Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, nil
	}

	var totalAmount float64
	query := fmt.Sprintf("SELECT SUM(p.price *c.quantity) FROM %s p JOIN %s c ON p.id =c.product_id WHERE user_id = $1", productsTable, cartsTable)
	row := r.db.QueryRow(c, query, userId)
	if err := row.Scan(&totalAmount); err != nil {
		return 0, err
	}
	return totalAmount, nil
}

func (r *UserRepository) Minus(c context.Context, userId uint, productId int) error {
	row := fmt.Sprintf("SELECT quantity from %s  WHERE user_id =$1 AND product_id =$2",
		cartsTable)
	var quantity int
	queryRow := r.db.QueryRow(c, row, userId, productId)
	err := queryRow.Scan(&quantity)
	if err != nil {
		return err
	}

	logrus.Print(quantity)
	if err != nil {
		return err
	}

	if quantity == 1 {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id =$1 AND product_id =$2",
			cartsTable)

		_, err = r.db.Exec(c, query, userId, productId)
		if err != nil {
			return err
		}
		return nil
	}

	query := fmt.Sprintf("UPDATE %s SET quantity = quantity - 1 WHERE user_id =$1 AND product_id =$2",
		cartsTable)

	_, err = r.db.Exec(c, query, userId, productId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Plus(c context.Context, userId uint, productId int) error {
	query := fmt.Sprintf("UPDATE %s SET quantity = quantity + 1 WHERE user_id =$1 AND product_id =$2",
		cartsTable)

	_, err := r.db.Exec(c, query, userId, productId)
	if err != nil {
		return err
	}
	return nil
}
