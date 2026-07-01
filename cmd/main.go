package main

import (
	"fmt"
	"net/http"

	"github.com/bekontaii/Online-Shop-Go/internal/cart"
	"github.com/bekontaii/Online-Shop-Go/internal/category"
	"github.com/bekontaii/Online-Shop-Go/internal/order"
	"github.com/bekontaii/Online-Shop-Go/internal/product"
	"github.com/bekontaii/Online-Shop-Go/internal/user"
	"github.com/bekontaii/Online-Shop-Go/pkg/database"
)

func main() {
	db := database.NewPostgresDB()

	// Users
	userRepo := user.NewPostgresRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Cart
	cartRepo := cart.NewPostgresRepository(db)
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)

	// Categories
	categoryRepo := category.NewPostgresRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// Products
	productRepo := product.NewPostgresRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Orders (depends on cartService)
	orderRepo := order.NewPostgresRepository(db)
	orderService := order.NewService(orderRepo, cartService)
	orderHandler := order.NewHandler(orderService)

	// HTTP Multiplexer
	mux := http.NewServeMux()

	// Register Routes
	userHandler.RegisterRoutes(mux)
	cartHandler.CartHandler(mux)
	categoryHandler.RegisterRoutes(mux)
	productHandler.RegisterRoutes(mux)
	orderHandler.RegisterRoutes(mux)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
