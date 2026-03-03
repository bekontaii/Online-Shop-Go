package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}
func (s *Service) Register(user *User) (*User, error) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return nil, errors.New("username and password are required")
	}
	usr, err := s.repo.GetUserByEmail(user.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}
	if err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	createdUser.Password = ""
	return createdUser, nil
}
