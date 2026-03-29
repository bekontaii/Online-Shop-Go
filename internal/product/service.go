package product

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("forbidden")
var ErrInvalidInput = errors.New("invalid input")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) CreateProduct(ctx context.Context, userID int64, role string, req CreateProductRequest) (int64, error) {
	if role != "admin" && role != "seller" {
		return 0, ErrForbidden
	}
	if req.Name == "" {
		return 0, ErrInvalidInput
	}
	if req.Price <= 0 {
		return 0, ErrInvalidInput
	}
	if req.Stock < 0 {
		return 0, ErrInvalidInput
	}
	product := Product{
		OwnerID:     userID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	id, err := s.repo.Create(ctx, product)
	if err != nil {
		return 0, err
	}
	return id, nil
}
