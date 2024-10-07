package likes

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"
	"y_net/internal/services/users"

	"github.com/google/uuid"
)

type likeRepositoryI interface {
	likePost(userId uuid.UUID, postId uuid.UUID) error
	unlikePost(userId uuid.UUID, postId uuid.UUID) error
	getFromPost(postId uuid.UUID) ([]users.User, error)
}

type likeRepository struct{}

func (i *likeRepository) likePost(userId uuid.UUID, postId uuid.UUID) error {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO likes (user_id, post_id) VALUES ($1, $2)",
		userId, postId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}

	return nil
}

func (i *likeRepository) unlikePost(userId uuid.UUID, postId uuid.UUID) error {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(context.Background(), "DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userId, postId)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	return nil
}

func (i *likeRepository) getFromPost(postId uuid.UUID) ([]users.User, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	rows, err := tx.Query(context.Background(), "SELECT u.id, u.username, u.full_name, u.avatar FROM likes l JOIN users u ON l.user_id = u.id WHERE l.post_id = $1", postId)
	if err != nil {
		return nil, fmt.Errorf("failed to select likes: %w", err)
	}
	defer rows.Close()

	var userLikes []users.User
	for rows.Next() {
		var user users.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userLikes = append(userLikes, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return userLikes, nil
}
