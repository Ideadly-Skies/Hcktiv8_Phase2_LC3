package handler

import (
	"time"
	"github.com/labstack/echo/v4"
	config "w4/lc3/config/database"
	utils "w4/lc3/utils"
	"net/http"
	"context"
)

type Order struct {
	OrderID    int       `json:"order_id"`
	UserID     int       `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type AddOrderResponse struct {
	Message   string  `json:"message"`
	OrderID   int     `json:"order_id"`
	TotalPrice float64 `json:"total_price"`
}

// @Summary Get User Orders
// @Description Retrieve a list of all orders for the logged-in user.
// @Tags Orders
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of user orders"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/orders [get]
func GetOrders(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	// Query to fetch user orders
	query := "SELECT order_id, user_id, total_price, created_at FROM orders WHERE user_id = $1"
	rows, err := config.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve orders"})
	}
	defer rows.Close()

	// Parse the result
	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.OrderID, &order.UserID, &order.TotalPrice, &order.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error scanning orders"})
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"orders": orders})
}

// @Summary Add a New Order
// @Description Place a new order for the logged-in user. Cart items are processed, and the cart is cleared after order creation.
// @Tags Orders
// @Accept  json
// @Produce  json
// @Success 201 {object} AddOrderResponse "Order placed successfully"
// @Failure 400 {object} map[string]string "Bad Request - Cart is empty"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/orders [post]
func AddOrder(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := utils.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	// Step 1: Fetch all cart items for the user
	queryCart := "SELECT product_id, quantity FROM carts WHERE user_id = $1"
	rows, err := config.Pool.Query(context.Background(), queryCart, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch cart items"})
	}
	defer rows.Close()

	var cartItems []CartItem
	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error parsing cart data"})
		}
		cartItems = append(cartItems, item)
	}

	// Check if cart is empty
	if len(cartItems) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Cart is empty"})
	}

	// Step 2: Calculate total price (mock product price as $10 for demonstration)
	var totalPrice float64
	for _, item := range cartItems {
		productPrice := 10.0 // Replace this with a query to fetch actual product prices
		totalPrice += float64(item.Quantity) * productPrice
	}

	// Step 3: Insert new order into the orders table
	queryOrder := "INSERT INTO orders (user_id, total_price) VALUES ($1, $2) RETURNING order_id"
	var orderID int
	err = config.Pool.QueryRow(context.Background(), queryOrder, userID, totalPrice).Scan(&orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create order"})
	}

	// Step 4: Clear the user's cart
	queryDeleteCart := "DELETE FROM carts WHERE user_id = $1"
	_, err = config.Pool.Exec(context.Background(), queryDeleteCart, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to clear cart"})
	}

	// Step 5: Return success response
	return c.JSON(http.StatusCreated, AddOrderResponse{
		Message:    "Order placed successfully",
		OrderID:    orderID,
		TotalPrice: totalPrice,
	})
}