package product

import "time"

type Product struct {
	ID          int64
	OwnerID     int64
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
