package order

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bekontaii/Online-Shop-Go/internal/cart"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateOrderWithTx(ctx context.Context, userID int64, cartItems []cart.CartItem) (*Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalPrice float64
	type productStockPrice struct {
		id    int64
		price float64
		qty   int
	}
	var itemsToInsert []productStockPrice

	for _, item := range cartItems {
		queryProduct := `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`
		var stock int
		var price float64
		err := tx.QueryRowContext(ctx, queryProduct, item.ProductID).Scan(&stock, &price)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product %d: %w", item.ProductID, err)
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product ID %d (have %d, request %d)", item.ProductID, stock, item.Quantity)
		}

		totalPrice += price * float64(item.Quantity)

		itemsToInsert = append(itemsToInsert, productStockPrice{
			id:    int64(item.ProductID),
			price: price,
			qty:   item.Quantity,
		})
	}

	queryOrder := `INSERT INTO orders (user_id, total_price, status) VALUES ($1, $2, 'pending') RETURNING id, created_at, updated_at`
	var order Order
	order.UserID = userID
	order.TotalPrice = totalPrice
	order.Status = "pending"
	err = tx.QueryRowContext(ctx, queryOrder, userID, totalPrice).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range itemsToInsert {
		queryDecr := `UPDATE products SET stock = stock - $1 WHERE id = $2`
		_, err := tx.ExecContext(ctx, queryDecr, item.qty, item.id)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for product %d: %w", item.id, err)
		}

		queryOrderItem := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`
		var orderItemID int64
		err = tx.QueryRowContext(ctx, queryOrderItem, order.ID, item.id, item.qty, item.price).Scan(&orderItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item for product %d: %w", item.id, err)
		}

		order.Items = append(order.Items, OrderItem{
			ID:        orderItemID,
			OrderID:   order.ID,
			ProductID: item.id,
			Quantity:  item.qty,
			Price:     item.price,
		})
	}

	queryClearCart := `DELETE FROM cart_items WHERE user_id = $1`
	_, err = tx.ExecContext(ctx, queryClearCart, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*Order, error) {
	queryOrder := `SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE id = $1`
	var o Order
	err := r.db.QueryRowContext(ctx, queryOrder, id).Scan(
		&o.ID,
		&o.UserID,
		&o.TotalPrice,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryItems := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, queryItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *PostgresRepository) ListByUserID(ctx context.Context, userID int64) ([]Order, error) {
	query := `SELECT id, user_id, total_price, status, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.TotalPrice,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]Order, error) {
	query := `SELECT id, user_id, total_price, status, created_at, updated_at FROM orders ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.TotalPrice,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
