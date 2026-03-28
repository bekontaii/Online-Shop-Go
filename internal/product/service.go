package product

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) CreateProduct(ctx context.Context, UserID int64, role string, req string) (int64, error) {
	return s.repo.CreateProduct(ctx, Product{})
}
