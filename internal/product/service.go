package product

import (
	"context"
	"errors"
)

var ErrForbidden = errors.New("forbidden")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) CreateProduct(ctx context.Context, userID int64, role string, req CreateProductRequest) (int64, error) {
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
