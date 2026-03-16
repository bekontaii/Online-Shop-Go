package cart

import "database/sql"

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
