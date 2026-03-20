package cart

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetCart(userID int) ([]CartItem, error) {
	return s.repo.GetCartByUserID(userID)
}
func (s *Service) AddToCart(ctx context.Context, userID int, productID int, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return s.repo.UpsertCartItem(ctx, userID, productID, quantity)
}
