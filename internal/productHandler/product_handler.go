package handler

import (
	"github.com/labstack/echo/v4"
	config "w4/lc3/config/database"
	// utils "w4/lc3/utils"
	"net/http"
	"context"
)

type Product struct {
	ProductID   int     `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// @Summary Get All Products
// @Description Retrieve a list of all available products.
// @Tags Products
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of products retrieved successfully"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products [get]
func GetAllProducts(c echo.Context) error {
	// Query to fetch all products
	query := "SELECT product_id, name, description, price FROM products"

	rows, err := config.Pool.Query(context.Background(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve products"})
	}
	defer rows.Close()

	// Parse the results into a list of products
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductID, &product.Name, &product.Description, &product.Price); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error scanning product data"})
		}
		products = append(products, product)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"products": products})
}

// @Summary Get Product by ID
// @Description Retrieve a product's details by its ID.
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} Product "Product retrieved successfully"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products/{id} [get]
func GetProductByID(c echo.Context) error {
	// Extract product ID from URL params
	productID := c.Param("id")

	// Query to fetch product by ID
	query := "SELECT product_id, name, description, price FROM products WHERE product_id = $1"

	var product Product
	err := config.Pool.QueryRow(context.Background(), query, productID).
		Scan(&product.ProductID, &product.Name, &product.Description, &product.Price)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	return c.JSON(http.StatusOK, product)
}