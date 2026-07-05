package category

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrCategoryNotFound  = errors.New("category not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateCategory(ctx context.Context, role string, req CreateCategoryRequest) (int64, error) {
	if role != "admin" {
		return 0, ErrForbidden
	}
	if req.Name == "" {
		return 0, ErrInvalidInput
	}
	cat := Category{
		Name:        req.Name,
		Description: req.Description,
	}
	return s.repo.Create(ctx, cat)
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return cat, nil
}

func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.List(ctx)
}

func (s *Service) UpdateCategory(ctx context.Context, role string, id int64, req UpdateCategoryRequest) error {
	if role != "admin" {
		return ErrForbidden
	}

	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	if req.Name != nil {
		if *req.Name == "" {
			return ErrInvalidInput
		}
		cat.Name = *req.Name
	}
	if req.Description != nil {
		cat.Description = *req.Description
	}

	return s.repo.Update(ctx, cat)
}

func (s *Service) DeleteCategory(ctx context.Context, role string, id int64) error {
	if role != "admin" {
		return ErrForbidden
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}
