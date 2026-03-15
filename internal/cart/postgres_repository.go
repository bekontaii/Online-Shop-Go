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
	rows, err := r.db.Query("SELECT * FROM cart WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CartItem
}
