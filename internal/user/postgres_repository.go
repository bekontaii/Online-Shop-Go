package user

import (
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (repo *PostgresRepository) CreateUser(user *User) (*User, error) {

	query := `
	INSERT INTO users (username, email, password, role)
	VALUES ($1,$2,$3,$4)
	RETURNING id, username, email, role
	`

	err := repo.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *PostgresRepository) GetUserByEmail(email string) (*User, error) {

	var user User

	query := `
	SELECT id,username,email,password,role
	FROM users
	WHERE email=$1
	`

	err := repo.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
