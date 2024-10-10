package likes

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type likeRepository interface {
	likePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	unlikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	userLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	getFromPost(ctx context.Context, postId uuid.UUID) ([]shared.User, error)
}

type likeRepositoryImpl struct{}

func (i *likeRepositoryImpl) likePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
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

	_, err = tx.Exec(
		ctx,
		"INSERT INTO likes (user_id, post_id) VALUES ($1, $2)",
		userId, postId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}

	return nil
}

func (i *likeRepositoryImpl) unlikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
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

	_, err = tx.Exec(ctx, "DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userId, postId)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	return nil
}

func (i *likeRepositoryImpl) userLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)", userId, postId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user liked post: %w", err)
	}

	return exists, nil
}

func (i *likeRepositoryImpl) getFromPost(ctx context.Context, postId uuid.UUID) ([]shared.User, error) {
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

	rows, err := tx.Query(ctx, "SELECT u.id, u.username, u.full_name, u.avatar FROM likes l JOIN users u ON l.user_id = u.id WHERE l.post_id = $1", postId)
	if err != nil {
		return nil, fmt.Errorf("failed to select likes: %w", err)
	}
	defer rows.Close()

	var userLikes []shared.User
	for rows.Next() {
		var user shared.User
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
