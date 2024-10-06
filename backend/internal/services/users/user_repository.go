package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	database "y_net/internal/database/postgres"
)

type userRepositoryI interface {
	create(user User) (uuid.UUID, error)
	get(idList []uuid.UUID) ([]User, error)
	getBySearch(searchStr string) ([]User, error)
	update(user User, id uuid.UUID) error
	delete(id uuid.UUID) error
}

type userRepository struct{}

func (i *userRepository) create(user User) (uuid.UUID, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var id uuid.UUID
	err = tx.QueryRow(
		context.Background(),
		"INSERT INTO users (username, password, email, full_name, description, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

func (i *userRepository) get(idList []uuid.UUID) ([]User, error) {
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

func (i *userRepository) getBySearch(searchStr string) ([]User, error) {
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

func (i *userRepository) update(user User, id uuid.UUID) error {
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
			user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar, id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		_, err = tx.Exec(
			context.Background(),
			"UPDATE users SET username = $1, email = $2, full_name = $3, description = $4, avatar = $5 WHERE id = $6",
			user.Username, user.Email, user.FullName, user.Description, user.Avatar, id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

func (i *userRepository) delete(id uuid.UUID) error {
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
