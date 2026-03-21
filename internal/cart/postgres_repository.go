package cart

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
func (r *PostgresRepository) GetCartByUserID(userID int) ([]CartItem, error) {

	query := `
		SELECT product_id, quantity
		FROM cart_items
		WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem

	for rows.Next() {
		var item CartItem

		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
func (r *PostgresRepository) UpsertCartItem(ctx context.Context, userID int, productID int, quantity int) error {
	query := `INSERT INTO cart (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id)
DO UPDATE SET quantity = cart.quantity + EXCLUDED.quantity;`
	_, err := r.db.ExecContext(ctx, query, userID, productID, quantity)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostgresRepository) RemoveCartItem(ctx context.Context, userID int, productID int) error {
}
