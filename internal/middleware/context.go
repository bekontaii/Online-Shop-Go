package middleware

import "net/http"

func GetUserID(r *http.Request) (int64, bool) {
	value := r.Context().Value(UserIDKey)
	if value == nil {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}
