package product

import "time"

type Product struct {
	ID          int
	Name        string
	Price       float64
	Description string
	CreatedDate time.Time
	UpdatedDate time.Time
}
