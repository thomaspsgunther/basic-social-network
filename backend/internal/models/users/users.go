package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	database "y_net/internal/database/postgres"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password,omitempty"`
	Email         *string   `json:"email" db:"email" unique:"true"`
	FullName      *string   `json:"full_name" db:"full_name"`
	Description   *string   `json:"description" db:"description"`
	Avatar        *string   `json:"avatar" db:"avatar"`
	FollowerCount int       `json:"follower_count" db:"follower_count"`
}

func (user *User) Create() error {
	if user.Username == "" || user.Password == "" {
		return fmt.Errorf("username and password must not be empty")
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO users (username, password, email, full_name, description, avatar) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func GetAll() ([]User, error) {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	rows, err := tx.Query(context.Background(), "SELECT id, username, full_name, avatar FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return users, nil
}

func Get(id uuid.UUID) (User, error) {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return User{}, err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return User{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var user User
	err = tx.QueryRow(
		context.Background(),
		"SELECT id, username, email, full_name, description, avatar, follower_count FROM users WHERE id = $1",
		id).Scan(&user.ID, &user.Username)
	if err != nil {
		return User{}, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}

func Update(user User, id uuid.UUID) error {
	user.ID = id

	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	if user.Password != "" {
		_, err = tx.Exec(
			context.Background(),
			"UPDATE users SET username = $1, password = $2, email = $3, full_name = $4, description = $5, avatar = $6 WHERE id = $7",
			user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar, user.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		_, err = tx.Exec(
			context.Background(),
			"UPDATE users SET username = $1, email = $2, full_name = $3, description = $4, avatar = $5 WHERE id = $6",
			user.Username, user.Email, user.FullName, user.Description, user.Avatar, user.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

func Delete(id uuid.UUID) error {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (user *User) Authenticate() (bool, error) {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var hashedPassword string
	err = tx.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", user.Username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	return CheckPasswordHash(user.Password, hashedPassword), nil
}

// GetUsernameByUserID checks if a user exists in database by given id
func GetUsernameByUserID(id uuid.UUID) (string, error) {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var username string
	err = tx.QueryRow(context.Background(), "SELECT username FROM users WHERE id = $1", id).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// GetUserIdByUsername checks if a user exists in database by given username
func GetUserIdByUsername(username string) (uuid.UUID, error) {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return uuid.UUID{}, err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var id uuid.UUID
	err = tx.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.UUID{}, &WrongUsernameOrPasswordError{}
		} else {
			return uuid.UUID{}, err
		}
	}

	return id, nil
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares raw password with its hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
