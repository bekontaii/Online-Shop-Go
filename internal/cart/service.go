package cart

import "errors"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetCart(userID int) ([]CartItem, error) {
	return s.repo.GetCartByUserID(userID)
}
func (s *Service) AddToCart(userID int, productID int, quantity int) error {
	if quantity < 0 {
		return errors.New("quantity must be greater than 0")
	}
	return s.repo.UpsertCartItem(userID, productID, quantity)
}
