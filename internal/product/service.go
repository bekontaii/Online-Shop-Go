package product

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidInput     = errors.New("invalid input")
	ErrProductNotFound = errors.New("product not found")
)

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
		CategoryID:  req.CategoryID,
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

func (s *Service) GetProduct(ctx context.Context, id int64) (*Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) ListProducts(ctx context.Context, categoryID *int64, search string) ([]Product, error) {
	return s.repo.List(ctx, categoryID, search)
}

func (s *Service) UpdateProduct(ctx context.Context, userID int64, role string, productID int64, req UpdateProductRequest) error {
	p, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	// Permission check: only admin or product owner can update
	if role != "admin" && p.OwnerID != userID {
		return ErrForbidden
	}

	// Update fields if provided
	if req.Name != nil {
		if *req.Name == "" {
			return ErrInvalidInput
		}
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return ErrInvalidInput
		}
		p.Price = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return ErrInvalidInput
		}
		p.Stock = *req.Stock
	}
	if req.CategoryID != nil {
		p.CategoryID = *req.CategoryID
	}

	return s.repo.Update(ctx, p)
}

func (s *Service) DeleteProduct(ctx context.Context, userID int64, role string, productID int64) error {
	p, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	// Permission check: only admin or product owner can delete
	if role != "admin" && p.OwnerID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, productID)
}
