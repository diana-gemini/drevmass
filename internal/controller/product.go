package controller

import (
	"net/http"
	"strconv"

	"github.com/diana-gemini/drevmass/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	ProductRepository models.ProductRepository
}

// Create Product godoc
// @Summary Create product
// @Security ApiKeyAuth
// @Tags admin-product-controller
// @ID create-product
// @Accept json
// @Produce json
// @Param name formData string true "product name"
// @Param image formData file true "poster"
// @Param price formData string true "price"
// @Param height formData string true "height"
// @Param size formData string true "size"
// @Param instruction formData string false "instruction"
// @Param description formData string true "description"
// @Param recommended_products formData []string false "recommended products"
// @Success 200 {integer} integer
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/products [post]
func (h *ProductController) createProduct(c *gin.Context) {

	var input models.Product

	name := c.PostForm("name")
	image, err := c.FormFile("image")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid image param")
		return
	}
	price := c.PostForm("price")
	pricee, err := strconv.ParseFloat(price, 64)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid price param")
		return
	}
	path := "files//product_images//" + image.Filename
	c.SaveUploadedFile(image, path)
	height := c.PostForm("height")
	size := c.PostForm("size")
	instruction := c.PostForm("instruction")
	description := c.PostForm("description")
	recommendedProducts := c.PostFormArray("recommended_products")

	logrus.Print(recommendedProducts)

	input = models.Product{
		Name:                name,
		Image:               path,
		Price:               pricee,
		Height:              height,
		Size:                size,
		Instruction:         instruction,
		Description:         description,
		RecommendedProducts: recommendedProducts}

	id, err := h.ProductRepository.Create(c, input)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to create the product")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary All products
// @Security ApiKeyAuth
// @Tags admin-product-cintroller
// @Accept json
// @Produce json
// @Success 200 {object} models.GetProducts
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [get]
func (h *ProductController) getAlProdacts(c *gin.Context) {

	products, err := h.ProductRepository.GetAll(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to get products")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Data": products,
	})

}

// @Summary edit the product
// @Security ApiKeyAuth
// @Tags admin-product-controller
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Param name formData string false "product name"
// @Param image formData file false "poster"
// @Param price formData string false "price"
// @Param height formData string false "height"
// @Param size formData string false "size"
// @Param instruction formData string false "instruction"
// @Param description formData string false "description"
// @Param recommended_products formData []string false "recommended products"
// @Success 200 {integer} integer
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/products/{id} [put]
func (h *ProductController) updateProduct(c *gin.Context) {

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	name := c.PostForm("name")

	image, err := c.FormFile("image")
	var path string
	if err != nil {
		product, _, err := h.ProductRepository.GetById(c, productId)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, "invalid image param")
			return
		}
		path = product.Image
	} else {
		path = "files//product_images//" + image.Filename
		c.SaveUploadedFile(image, path)

	}
	price := c.PostForm("price")
	var pricee float64
	if price != "" {
		priceee, err := strconv.ParseFloat(price, 64)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "invalid price param")
			return
		}
		pricee = priceee

	}

	height := c.PostForm("height")
	size := c.PostForm("size")
	instruction := c.PostForm("instruction")
	description := c.PostForm("description")
	recommendedProducts := c.PostFormArray("recommended_products")

	input := models.UpdateProduct{
		Name:                name,
		Image:               path,
		Price:               pricee,
		Height:              height,
		Size:                size,
		Instruction:         instruction,
		Description:         description,
		RecommendedProducts: recommendedProducts}

	logrus.Print(input.Instruction)

	if err := h.ProductRepository.Update(c, productId, input); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to edit the product")
		return
	}

	c.JSON(http.StatusOK, models.StatusResponse{
		Status: "ok",
	})
}

// @Summary delete product
// @Security ApiKeyAuth
// @Tags admin-product-controller
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Success 200 {object} models.StatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/products/{id} [delete]
func (h *ProductController) deleteProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	if err := h.ProductRepository.Delete(c, productId); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, "failed to delete the product")
		return
	}
	c.JSON(http.StatusOK, models.StatusResponse{
		Status: "ok",
	})
}

// @Summary product by id
// @Security ApiKeyAuth
// @Tags admin-product-controller
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.GetProducts
// @Failure 500 {object} models.GetProducts
// @Router /products/{id} [get]
func (h *ProductController) getProdactById(c *gin.Context) {

	prodactId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	product, products, err := h.ProductRepository.GetById(c, prodactId)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Product":              product,
		"Recommended products": products,
	})

}
