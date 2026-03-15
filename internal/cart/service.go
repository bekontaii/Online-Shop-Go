package cart

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetCart(userID int) ([]CartItem, error) {
	return s.repo.GetCartByUserID(userID)
}
