package users

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseI interface {
	Create(user User) (uuid.UUID, error)
	Get(idList []uuid.UUID) ([]User, error)
	GetBySearch(searchStr string) ([]User, error)
	Update(user User, id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type userUsecase struct {
	usecase    UserUsecaseI
	repository userRepositoryI
}

func NewUserUsecase() UserUsecaseI {
	return &userUsecase{
		usecase:    &userUsecase{},
		repository: &userRepository{},
	}
}

func (i *userUsecase) Create(user User) (uuid.UUID, error) {
	if user.Username == "" || user.Password == "" {
		return uuid.Nil, fmt.Errorf("username and password must not be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	id, err := i.repository.create(user)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (i *userUsecase) Get(idList []uuid.UUID) ([]User, error) {
	users, err := i.repository.get(idList)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *userUsecase) GetBySearch(searchStr string) ([]User, error) {
	users, err := i.repository.getBySearch(searchStr)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *userUsecase) Update(user User, id uuid.UUID) error {
	if user.Username == "" {
		return fmt.Errorf("username must not be empty")
	}

	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	err := i.repository.update(user, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *userUsecase) Delete(id uuid.UUID) error {
	err := i.repository.delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) Authenticate() (bool, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var hashedPassword string
	err = tx.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", user.Username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	return checkPasswordHash(user.Password, hashedPassword), nil
}

// GetUsernameByUserID checks if a user exists in database by given id
func GetUsernameByUserID(id uuid.UUID) (string, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var username string
	err = tx.QueryRow(context.Background(), "SELECT username FROM users WHERE id = $1", id).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// GetUserIdByUsername checks if a user exists in database by given username
func GetUserIdByUsername(username string) (uuid.UUID, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var id uuid.UUID
	err = tx.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, &WrongUsernameOrPasswordError{}
		} else {
			return uuid.Nil, err
		}
	}

	return id, nil
}

// HashPassword hashes given password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares raw password with its hashed values
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
