package product

import "context"

type Repository interface {
	Create(ctx context.Context, product Product) (int64, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, categoryID *int64, search string) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int64) error
}
