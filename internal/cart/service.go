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
func (s *Service) RemoveFromCart(ctx context.Context, userID int, productID int) error {
	if userID <= 0 {
		return errors.New("user_id must be greater than 0")
	}
	if productID <= 0 {
		return errors.New("product_id must be greater than 0")
	}
	return s.repo.RemoveCartItem(ctx, userID, productID)
}
func (s *Service) UpdateCartItem(ctx context.Context, userID int, productID int, quantity int) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if productID <= 0 {
		return ErrInvalidProductID
	}
	if quantity < 0 {
		return ErrInvalidQuantity
	}

	if quantity == 0 {
		return s.repo.RemoveCartItem(ctx, userID, productID)
	}

	return s.repo.UpdateCartItem(ctx, userID, productID, quantity)
}
