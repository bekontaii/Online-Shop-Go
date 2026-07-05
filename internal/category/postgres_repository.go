package category

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

func (r *PostgresRepository) Create(ctx context.Context, category Category) (int64, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, category.Name, category.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	var cat Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Description,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *PostgresRepository) Update(ctx context.Context, category *Category) error {
	query := `UPDATE categories SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.Description, category.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
