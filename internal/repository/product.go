package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) models.ProductRepository {
	return &ProductRepository{db: db}
}

const (
	productsTable            = "products"
	recommendedProductsTable = "recommended_products"
)

func (r *ProductRepository) Create(c context.Context, product models.Product) (int, error) {
	var recommendedProductss []string
	for _, v := range product.RecommendedProducts {
		recommendedProductss = strings.Split(v, ",")

	}

	var recommendedProductsArray []int
	for _, v := range recommendedProductss {
		pId, err := strconv.Atoi(v)
		if err != nil {

			return 0, err
		}
		recommendedProductsArray = append(recommendedProductsArray, pId)
	}

	var id int
	productQuery := fmt.Sprintf("INSERT INTO %s (name, image, price,height, size, instruction, description) values($1, $2,$3,$4,$5,$6,$7) RETURNING id", productsTable)
	row := r.db.QueryRow(c, productQuery, product.Name, product.Image, product.Price, product.Height, product.Size, product.Instruction, product.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	logrus.Print(recommendedProductsArray)
	for _, rProduct := range recommendedProductsArray {

		recommendedQuery := fmt.Sprintf("INSERT INTO %s (product_id, recommended_product) values ($1,$2)", recommendedProductsTable)
		_, err := r.db.Exec(c, recommendedQuery, id, rProduct)
		if err != nil {
			return id, err
		}
	}
	return id, nil

}

func (r *ProductRepository) GetAll(c context.Context) ([]models.GetProducts, error) {

	query := fmt.Sprintf("SELECT id, name, image, price FROM %s order by id", productsTable)
	logrus.Print(query)
	rows, err := r.db.Query(c, query)
	if err != nil {
		return nil, err
	}

	var products = []models.GetProducts{}
	for rows.Next() {
		var product models.GetProducts
		err := rows.Scan(&product.Id, &product.Name, &product.Image, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, err
}

func (r *ProductRepository) Update(c context.Context, productId int, input models.UpdateProduct) error {
	var recommendedProductss []string
	for _, v := range input.RecommendedProducts {
		recommendedProductss = strings.Split(v, ",")

	}

	var recommendedProductsArray []int
	for _, v := range recommendedProductss {
		pId, err := strconv.Atoi(v)
		if err != nil {

			return err
		}
		recommendedProductsArray = append(recommendedProductsArray, pId)
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, input.Name)
		argId++
	}

	if input.Image != "" {
		setValues = append(setValues, fmt.Sprintf("image = $%d", argId))
		args = append(args, input.Image)
		argId++
	}

	if input.Price != 0 {
		setValues = append(setValues, fmt.Sprintf("price = $%d", argId))
		args = append(args, input.Price)
		argId++
	}

	if input.Height != "" {
		setValues = append(setValues, fmt.Sprintf("height = $%d", argId))
		args = append(args, input.Height)
		argId++
	}

	if input.Size != "" {
		setValues = append(setValues, fmt.Sprintf("size = $%d", argId))
		args = append(args, input.Size)
		argId++
	}

	if input.Instruction != "" {
		setValues = append(setValues, fmt.Sprintf("instruction = $%d", argId))
		args = append(args, input.Instruction)
		argId++
	}

	if input.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, input.Description)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id =$%d ",
		productsTable, setQuery, argId)

	args = append(args, productId)

	logrus.Print(setQuery)
	logrus.Print(query)
	logrus.Print(args)

	_, err := r.db.Exec(c, query, args...)
	if err != nil {
		return err
	}

	if recommendedProductsArray != nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE product_id = $1", recommendedProductsTable)
		_, err := r.db.Exec(c, query, productId)
		if err != nil {
			return err
		}
		for _, rProduct := range recommendedProductsArray {
			recommendedQuery := fmt.Sprintf("INSERT INTO %s (product_id, recommended_product) values ($1,$2)", recommendedProductsTable)
			_, err = r.db.Exec(c, recommendedQuery, productId, rProduct)
			if err != nil {
				return err
			}
		}
	}
	return err

}

func (r *ProductRepository) Delete(c context.Context, productId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1",
		productsTable)
	_, err := r.db.Exec(c, query, productId)
	return err
}

func (r *ProductRepository) GetById(c context.Context, productId int) (models.GetProduct, []models.GetProducts, error) {
	var product models.GetProduct
	query := fmt.Sprintf("SELECT id, name, image, price, height, size, instruction,description FROM %s WHERE id = $1", productsTable)
	logrus.Print(query)
	row := r.db.QueryRow(c, query, productId)
	err := row.Scan(&product.Id, &product.Name, &product.Image, &product.Price, &product.Height, &product.Size, &product.Instruction, &product.Description)
	if err != nil {
		return product, nil, err
	}

	var products []models.GetProducts
	rQuery := fmt.Sprintf("SELECT p.id, p.name, p.image, p.price FROM %s p JOIN %s rp ON p.id = rp.recommended_product WHERE  rp.product_id=$1  ", productsTable, recommendedProductsTable)
	logrus.Print(rQuery)
	rows, err := r.db.Query(c, rQuery, productId)
	if err != nil {
		return product, nil, err
	}

	for rows.Next() {
		var productt models.GetProducts
		err := rows.Scan(&productt.Id, &productt.Name, &productt.Image, &productt.Price)
		if err != nil {
			return product, nil, err
		}
		products = append(products, productt)

	}

	return product, products, err
}
