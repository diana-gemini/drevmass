package models

import (
	"context"
)

type Product struct {
	Id                  int      `json:"-" db:"id"`
	Name                string   `json:"name" binding:"required"`
	Image               string   `json:"image" binding:"required"`
	Price               float64  `json:"price" binding:"required"`
	Height              string   `json:"height" binding:"required"`
	Size                string   `json:"size" binding:"required"`
	Instruction         string   `json:"instruction"`
	Description         string   `json:"description" binding:"required"`
	RecommendedProducts []string `json:"recommended_products"`
}

type ProductRepository interface {
	Create(c context.Context, product Product) (int, error)
	GetAll(c context.Context) ([]GetProducts, error)
	Update(c context.Context, productId int, input UpdateProduct) error
	Delete(c context.Context, productId int) error
	GetById(c context.Context, productId int) (GetProduct, []GetProducts, error)
}

type GetProducts struct {
	Id    int     `json:"id" db:"id"`
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Price float64 `json:"price"`
}

type UpdateProduct struct {
	Name                string   `json:"name"`
	Image               string   `json:"image"`
	Price               float64  `json:"price"`
	Height              string   `json:"height"`
	Size                string   `json:"size"`
	Instruction         string   `json:"instruction"`
	Description         string   `json:"description"`
	RecommendedProducts []string `json:"recommended_products"`
}

type GetProduct struct {
	Id          int     `json:"id" db:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Height      string  `json:"height"`
	Size        string  `json:"size"`
	Instruction string  `json:"instruction"`
	Description string  `json:"description"`
}
