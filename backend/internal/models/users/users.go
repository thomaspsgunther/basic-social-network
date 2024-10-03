package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	database "y_net/internal/database/postgres"
)

type TokenJson struct {
	Token string `json:"token"`
}

type Users struct {
	UserList []*User
}

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password,omitempty"`
	Email         *string   `json:"email"`
	FullName      *string   `json:"fullName"`
	Description   *string   `json:"description"`
	Avatar        *string   `json:"avatar"`
	FollowerCount int       `json:"followeCount"`
}

func (user *User) Create() (uuid.UUID, error) {
	if user.Username == "" || user.Password == "" {
		return uuid.UUID{}, fmt.Errorf("username and password must not be empty")
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return uuid.UUID{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var id uuid.UUID
	err = tx.QueryRow(
		context.Background(),
		"INSERT INTO users (username, password, email, full_name, description, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar,
	).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

func Get(idList []uuid.UUID) ([]User, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var placeholders []string
	for i := 0; i < len(idList); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(idList)+1))
	}

	placeholdersClause := strings.Join(placeholders, ", ")
	query := fmt.Sprintf("SELECT id, username, full_name, avatar FROM users WHERE id IN (%s)", placeholdersClause)

	rows, err := tx.Query(context.Background(), query, idList)
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

func GetBySearch(searchStr string) ([]User, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	query := "SELECT id, username, full_name, avatar FROM users WHERE username ILIKE $1 OR full_name ILIKE $1"

	searchPattern := "%" + searchStr + "%"

	rows, err := tx.Query(context.Background(), query, searchPattern)
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

func Update(user User, id uuid.UUID) error {
	user.ID = id

	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
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
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
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
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
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
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
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
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return uuid.UUID{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
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
