package user

import "errors"

type Repository interface {
	CreateUser(user *User) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

var ErrUserNotFound = errors.New("user not found")

type InMemoryRepository struct {
	users map[string]*User
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryRepository) CreateUser(user *User) (*User, error) {
	if _, ok := r.users[user.Email]; ok {
		return nil, errors.New("user already exists")
	}

	r.users[user.Email] = user
	return user, nil
}

func (r *InMemoryRepository) GetUserByEmail(email string) (*User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}
