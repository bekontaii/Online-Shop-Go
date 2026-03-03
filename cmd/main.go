package main

import (
	"fmt"
	"github.com/bekontaii/Online-Shop-Go.git/internal/user"
	"net/http"
)

func main() {
	fmt.Println("Pet Project Online Shop")
	repo := user.NewUserRepository()
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	user.RegisterRoutes(muhandler)
}
