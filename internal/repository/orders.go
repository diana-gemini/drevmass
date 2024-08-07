package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/diana-gemini/drevmass/internal/models"
)

const (
	ordersTable = "orders"
	ordersItems = "order_items"
)

func (r *UserRepository) CreateOrder(c context.Context, userId int, totalAmount float64, bonus float64, overall float64, products []models.GetProductsFromCart) (int, error) {
	var id int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf("INSERT INTO %s (user_id, total_amount,bonus,overall,created_date) values($1, $2,$3,$4,$5) RETURNING id", ordersTable)
	row := r.db.QueryRow(c, query, userId, totalAmount, bonus, overall, currentTime)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	for _, v := range products {
		var itemId int
		itemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, product_id,quantity, price) values($1, $2,$3,$4) RETURNING id", ordersItems)
		row := r.db.QueryRow(c, itemsQuery, id, v.Id, v.Quantity, v.Price)
		if err := row.Scan(&itemId); err != nil {
			return 0, err
		}
	}
	return id, nil
}

// GetAll implements models.OrderRepository.
func (r *UserRepository) GetAllOrders(c context.Context) ([]models.OrderForAdmin, error) {
	var orders []models.OrderForAdmin
	query := fmt.Sprintf(`select o.id, o.user_id, o.total_amount, o.bonus, o.overall, o.created_date, oi.product_id,p.name,p.image,p.price,oi.quantity 
						  FROM %s o JOIN %s oi ON o.id = oi.order_id  
						  JOIN %s p ON oi.product_id = p.id
						  ORDER BY o.id DESC`,
		ordersTable, ordersItems, productsTable)

	rows, err := r.db.Query(c, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.OrderForAdmin
		err := rows.Scan(&order.Id, &order.UserId, &order.TotalAmount, &order.Bonus, &order.Overall,
			&order.CreatedAt, &order.ProductId, &order.Name, &order.Image, &order.Price, &order.Quantity)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err != nil {
		return nil, err
	}

	return orders, err
}

func (r *UserRepository) GetOrderById(c context.Context, id int) ([]models.OrderForAdmin, error) {
	var orders []models.OrderForAdmin
	query := fmt.Sprintf(`select o.id, o.user_id, o.total_amount, o.bonus, o.overall, o.created_date, oi.product_id,p.name,p.image,p.price,oi.quantity 
						  FROM %s o JOIN %s oi ON o.id = oi.order_id  
						  JOIN %s p ON oi.product_id = p.id
						  WHERE o.id = %d	
						  `,
		ordersTable, ordersItems, productsTable, id)

	rows, err := r.db.Query(c, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.OrderForAdmin
		err := rows.Scan(&order.Id, &order.UserId, &order.TotalAmount, &order.Bonus, &order.Overall,
			&order.CreatedAt, &order.ProductId, &order.Name, &order.Image, &order.Price, &order.Quantity)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err != nil {
		return nil, err
	}

	return orders, err
}

// GetAllForUser implements models.OrderRepository.
func (*UserRepository) GetAllOrdersForUser(c context.Context, userId int) ([]models.OrderForUser, error) {
	panic("unimplemented")
}
