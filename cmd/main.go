package main

import (
	"fmt"
	"github.com/bekontaii/Online-Shop-Go/internal/cart"
	"github.com/bekontaii/Online-Shop-Go/internal/user"
	"github.com/bekontaii/Online-Shop-Go/pkg/database"
	"net/http"
)

func main() {
	db := database.NewPostgresDB()

	repo := user.NewPostgresRepository(db)

	service := user.NewService(repo)

	handler := user.NewHandler(service)

	mux := http.NewServeMux()

	handler.RegisterRoutes(mux)
	cartRepo := cart.NewPostgresRepository(db)
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)
	cartHandler.CartHandler(mux)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", mux)

}
