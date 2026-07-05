package category

import "context"

type Repository interface {
	Create(ctx context.Context, category Category) (int64, error)
	GetByID(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int64) error
}
