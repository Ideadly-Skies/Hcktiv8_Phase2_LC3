package main

import (
	config "w4/lc3/config/database"
	cust_middleware "w4/lc3/internal/middleware"
	user_handler "w4/lc3/internal/userHandler"
	cart_handler "w4/lc3/internal/cartHandler"
	order_handler "w4/lc3/internal/orderHandler"
	product_handler "w4/lc3/internal/productHandler"
	"github.com/swaggo/echo-swagger"
	_ "w4/lc3/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main(){
	// migrate data to supabase
	// config.MigrateData()

	// connect to db
	config.InitDB()
	defer config.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// public routes
	e.POST("users/register", user_handler.Register)	
	e.POST("users/login", user_handler.Login)
	
	// products
	e.GET("products", product_handler.GetAllProducts)
	e.GET("products/:id", product_handler.GetProductByID)

	// protected routes //
	// carts
	e.GET("users/carts", cart_handler.GetCart, cust_middleware.JWTMiddleware)              
	e.POST("users/carts", cart_handler.AddToCart, cust_middleware.JWTMiddleware)           
	e.DELETE("users/carts/:id", cart_handler.DeleteCartItem, cust_middleware.JWTMiddleware) 

	// orders
	e.GET("users/orders", order_handler.GetOrders, cust_middleware.JWTMiddleware)
	e.POST("users/orders", order_handler.AddOrder, cust_middleware.JWTMiddleware)

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// start the server at 8080
	e.Logger.Fatal(e.Start(":8080"))
}