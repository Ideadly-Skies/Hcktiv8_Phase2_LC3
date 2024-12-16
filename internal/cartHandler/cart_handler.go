package handler

import (
	"time"
	"net/http"
	config "w4/lc3/config/database"
	utils "w4/lc3/utils"

	"context"
	"github.com/labstack/echo/v4"
)

// Cart struct
type Cart struct {
	CartID    int       `json:"cart_id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// AddToCartRequest struct
type AddToCartRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

// @Summary Retrieve cart items for the logged-in user
// @Description Get all cart items belonging to the authenticated user
// @Tags Carts
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of cart items"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to retrieve cart data"
// @Router /users/carts [get]
func GetCart(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	// Query to get cart data
	query := "SELECT cart_id, user_id, product_id, quantity, created_at FROM carts WHERE user_id = $1"
	rows, err := config.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve cart data"})
	}
	defer rows.Close()

	// Fetch cart data
	var cartItems []Cart
	for rows.Next() {
		var cart Cart
		if err := rows.Scan(&cart.CartID, &cart.UserID, &cart.ProductID, &cart.Quantity, &cart.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error scanning cart data"})
		}
		cartItems = append(cartItems, cart)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"cart": cartItems})
}

// @Summary Add an item to the user's cart
// @Description Add a product to the cart for the authenticated user
// @Tags Carts
// @Accept  json
// @Produce  json
// @Param request body AddToCartRequest true "Request body for adding a product to the cart"
// @Success 201 {object} map[string]string "Item added to cart"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to add to cart"
// @Router /users/carts [post]
func AddToCart(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	// Parse request body
	var req AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Insert into cart
	query := "INSERT INTO carts (user_id, product_id, quantity) VALUES ($1, $2, $3)"
	_, err = config.Pool.Exec(context.Background(), query, userID, req.ProductID, req.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add to cart"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Item added to cart"})
}

// @Summary Delete a specific item from the user's cart
// @Description Delete a cart item based on the cart ID for the authenticated user
// @Tags Carts
// @Accept  json
// @Produce  json
// @Param id path int true "Cart ID"
// @Success 200 {object} map[string]string "Item deleted from cart"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Cart item not found"
// @Failure 500 {object} map[string]string "Failed to delete item"
// @Router /users/carts/{id} [delete]
func DeleteCartItem(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	// Get cart ID from URL params
	cartID := c.Param("id")

	// Delete item where user_id and cart_id match
	query := "DELETE FROM carts WHERE cart_id = $1 AND user_id = $2"
	result, err := config.Pool.Exec(context.Background(), query, cartID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete item"})
	}

	// Check if any row was affected
	if result.RowsAffected() == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Cart item not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted from cart"})
}