package models

import "context"

type CartRepository interface {
	Add(c context.Context, userId uint, productId int) (int, error)
	Minus(c context.Context, userId uint, productId int) error
	Plus(c context.Context, userId uint, productId int) error
	Delete(c context.Context, userId uint) error
	GetAllFromCart(c context.Context, userId uint) ([]GetProductsFromCart, error)
	GetTotalAmout(c context.Context, userId uint) (float64, error)

	Create(c context.Context, userId int, totalAmount float64, bonus float64, products []GetProductsFromCart) (int, error)
	GetAllForUser(c context.Context, userId int) ([]OrderForUser, error)
	GetAll(c context.Context) ([]OrderForAdmin, error)
}

type GetProductsFromCart struct {
	Id       int     `json:"id" db:"id"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
