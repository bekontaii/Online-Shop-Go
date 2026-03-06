package auth

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type LoginRequest struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}
