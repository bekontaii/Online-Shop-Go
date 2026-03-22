package cart

import "context"

type Repository interface {
	GetCartByUserID(userID int) ([]CartItem, error)
	UpsertCartItem(ctx context.Context, userID int, productID int, quantity int) error
	RemoveCartItem(ctx context.Context, userID int, productID int) error
	UpdartCartItem(ctx context.Context, userID int, productID int, quantity int) error
}
