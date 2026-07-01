package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bekontaii/Online-Shop-Go/internal/cart"
)

var (
	ErrForbidden     = errors.New("forbidden")
	ErrOrderNotFound = errors.New("order not found")
	ErrEmptyCart     = errors.New("cart is empty")
)

type Service struct {
	repo        Repository
	cartService *cart.Service
}

func NewService(repo Repository, cartService *cart.Service) *Service {
	return &Service{
		repo:        repo,
		cartService: cartService,
	}
}

func (s *Service) Checkout(ctx context.Context, userID int64) (*Order, error) {
	cartItems, err := s.cartService.GetCart(int(userID))
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, ErrEmptyCart
	}

	return s.repo.CreateOrderWithTx(ctx, userID, cartItems)
}

func (s *Service) GetOrder(ctx context.Context, userID int64, role string, orderID int64) (*Order, error) {
	o, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Permission: Users can view only their own orders. Admins/Sellers can view any order.
	if role != "admin" && role != "seller" && o.UserID != userID {
		return nil, ErrForbidden
	}

	return o, nil
}

func (s *Service) ListOrders(ctx context.Context, userID int64, role string) ([]Order, error) {
	if role == "admin" || role == "seller" {
		return s.repo.ListAll(ctx)
	}
	return s.repo.ListByUserID(ctx, userID)
}

func (s *Service) UpdateStatus(ctx context.Context, role string, orderID int64, status string) error {
	// Only admin or seller can change order status
	if role != "admin" && role != "seller" {
		return ErrForbidden
	}

	if status == "" {
		return errors.New("status cannot be empty")
	}

	_, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderNotFound
		}
		return err
	}

	return s.repo.UpdateStatus(ctx, orderID, status)
}
