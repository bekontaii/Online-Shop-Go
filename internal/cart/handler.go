package cart

import (
	"fmt"
	"github.com/bekontaii/Online-Shop-Go/internal/middleware"
	"net/http"
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "%d", userID)
}
