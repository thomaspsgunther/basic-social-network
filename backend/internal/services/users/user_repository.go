package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	database "y_net/internal/database/postgres"
	"y_net/internal/services/shared"
)

type userRepository interface {
	create(ctx context.Context, user shared.User) (uuid.UUID, error)
	get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error)
	getBySearch(ctx context.Context, searchStr string) ([]shared.User, error)
	getPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error)
	update(ctx context.Context, user shared.User, id uuid.UUID) error
	delete(ctx context.Context, id uuid.UUID) error
}

type userRepositoryImpl struct{}

func (r *userRepositoryImpl) create(ctx context.Context, user shared.User) (uuid.UUID, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var id uuid.UUID
	err = tx.QueryRow(
		ctx,
		"INSERT INTO users (username, password, email, full_name, description, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

func (r *userRepositoryImpl) get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var placeholders []string
	for i := 0; i < len(idList); i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(idList)+1))
	}

	placeholdersClause := strings.Join(placeholders, ", ")
	query := fmt.Sprintf("SELECT id, username, full_name, avatar FROM users WHERE id IN (%s)", placeholdersClause)

	rows, err := tx.Query(ctx, query, idList)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	var users []shared.User
	for rows.Next() {
		var user shared.User
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

func (r *userRepositoryImpl) getBySearch(ctx context.Context, searchStr string) ([]shared.User, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	query := "SELECT id, username, full_name, avatar FROM users WHERE username ILIKE $1 OR full_name ILIKE $1"

	searchPattern := "%" + searchStr + "%"

	rows, err := tx.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	var users []shared.User
	for rows.Next() {
		var user shared.User
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

func (r *userRepositoryImpl) getPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		database.HandleTransaction(tx, err)
	}()

	query := `
        SELECT id, image
        FROM posts
        WHERE user_id = $1
        AND (created_at < $2 OR (created_at = $2 AND id < $3))
        ORDER BY created_at DESC, id DESC
        LIMIT $4
    `

	rows, err := tx.Query(ctx, query, userId, lastCreatedAt, lastId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	defer rows.Close()

	var posts []shared.Post
	for rows.Next() {
		var post shared.Post
		if err := rows.Scan(&post.ID, &post.Image); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return posts, nil
}

func (r *userRepositoryImpl) update(ctx context.Context, user shared.User, id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	if user.Password != "" {
		_, err = tx.Exec(
			ctx,
			"UPDATE users SET username = $1, password = $2, email = $3, full_name = $4, description = $5, avatar = $6 WHERE id = $7",
			user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar, id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		_, err = tx.Exec(
			ctx,
			"UPDATE users SET username = $1, email = $2, full_name = $3, description = $4, avatar = $5 WHERE id = $6",
			user.Username, user.Email, user.FullName, user.Description, user.Avatar, id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

func (r *userRepositoryImpl) delete(ctx context.Context, id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
