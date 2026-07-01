package order

import (
	"context"
	"github.com/bekontaii/Online-Shop-Go/internal/cart"
)

type Repository interface {
	CreateOrderWithTx(ctx context.Context, userID int64, cartItems []cart.CartItem) (*Order, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	ListByUserID(ctx context.Context, userID int64) ([]Order, error)
	ListAll(ctx context.Context) ([]Order, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}
