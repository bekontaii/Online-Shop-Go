package cart

import (
	"encoding/json"
	"errors"
	"github.com/bekontaii/Online-Shop-Go/internal/middleware"
	"net/http"
	"strconv"
)

var (
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidProductID = errors.New("invalid product id")
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	items, err := h.service.GetCart(int(userID))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}
func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	var req AddToCartRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if req.ProductID <= 0 || req.Quantity <= 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := h.service.AddToCart(r.Context(), int(userID), req.ProductID, req.Quantity)
	if err != nil {
		if errors.Is(err, ErrInvalidQuantity) ||
			errors.Is(err, ErrInvalidProductID) ||
			errors.Is(err, ErrInvalidUserID) {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = h.service.RemoveFromCart(r.Context(), int(userID), productID)
	if err != nil {
		if errors.Is(err, ErrInvalidProductID) ||
			errors.Is(err, ErrInvalidUserID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
func (h *Handler) CartHandler(mux *http.ServeMux) {
	cartHandler := middleware.JWTMiddleware

	mux.Handle("/api/cart", cartHandler(http.HandlerFunc(h.GetCart)))
	mux.Handle("/api/cart/add", cartHandler(http.HandlerFunc(h.AddToCart)))
	mux.Handle("/api/cart/delete", cartHandler(http.HandlerFunc(h.RemoveFromCart)))
}
