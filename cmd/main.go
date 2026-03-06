package main

import (
	"fmt"
	"github.com/bekontaii/Online-Shop-Go/internal/user"
	"net/http"
)

func main() {
	fmt.Println("Pet Project Online Shop")
	repo := user.NewInMemoryRepository()
	service := user.NewService(repo)
	handler := user.NewHandler(service)
	router := http.NewServeMux()
	handler.RegisterRoutes(router)
	http.ListenAndServe(":8080", router)
}
