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
	query := `INSERT INTO products (owner_id, category_id, name, description, price, stock)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`
	row := r.db.QueryRowContext(ctx, query,
		product.OwnerID,
		product.CategoryID,
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

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	query := `SELECT id, owner_id, category_id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`
	var p Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.OwnerID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresRepository) List(ctx context.Context, categoryID *int64, search string) ([]Product, error) {
	query := `SELECT id, owner_id, category_id, name, description, price, stock, created_at, updated_at 
FROM products 
WHERE ($1::BIGINT IS NULL OR category_id = $1)
  AND ($2::TEXT = '' OR name ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')
ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, categoryID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.ID,
			&p.OwnerID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *PostgresRepository) Update(ctx context.Context, product *Product) error {
	query := `UPDATE products 
SET category_id = $1, name = $2, description = $3, price = $4, stock = $5, updated_at = NOW() 
WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID,
	)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
