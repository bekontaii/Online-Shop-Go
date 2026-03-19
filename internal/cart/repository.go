package cart

type Repository interface {
	GetCartByUserID(userID int) ([]CartItem, error)
	UpsertCartItem(userID int, productID int, quantity int) error
}
