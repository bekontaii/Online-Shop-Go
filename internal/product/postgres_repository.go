package product

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}
func (r *PostgresRepository) Create(ctx context.Context, product Product) (int64, error) {
	query := `INSERT INTO products (owner_id, name, description, price, stock)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`
	row := r.db.QueryRowContext(ctx, query,
		product.OwnerID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
