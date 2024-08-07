package models

import "time"

type OrderForUser struct {
	Id          int       `json:"-" db:"id"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_date"`
	ProductId   int       `json:"product_id" db:"product_id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
}

type OrderWithProducts struct {
	Id          int                   `json:"id" db:"id"`
	TotalAmount float64               `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time             `json:"created_at" db:"created_date"`
	UserId      int                   `json:"user_id" db:"user_id"`
	Products    []GetProductsFromCart `json:"products"`
	Bonus       float64               `json:"bonus"`
	Overall     float64               `json:"overall"`
}

type OrderForAdmin struct {
	Id          int       `json:"-" db:"id"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_date"`
	UserId      int       `json:"user_id" db:"user_id"`
	ProductId   int       `json:"product_id" db:"product_id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Bonus       float64   `json:"bonus"`
	Overall     float64   `json:"overall"`
}
