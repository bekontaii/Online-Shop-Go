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
	_, err := s.repo.GetUserByEmail(user.Email)

	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	createdUser.Password = ""
	return createdUser, nil
}
func (s *Service) Login(email string, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("username or password is incorrect")
	}
	user.Password = ""
	return user, nil
}
