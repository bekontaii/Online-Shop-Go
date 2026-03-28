package product

import "context"

type Repository interface {
	CreateProduct(ctx context.Context, product Product) (int64, error)
}
