package cart

type Repository interface {
	GetCartByUserID(userID int) ([]CartItem, error)
}
