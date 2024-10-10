package followers

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type followerRepository interface {
	follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	userFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error)
	getFollowers(ctx context.Context, userId uuid.UUID) ([]shared.User, error)
	getFollowed(ctx context.Context, userId uuid.UUID) ([]shared.User, error)
}

type followerRepositoryImpl struct{}

func (r *followerRepositoryImpl) follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
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
		"INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2)",
		followerId, followedId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert follower: %w", err)
	}

	return nil
}

func (r *followerRepositoryImpl) unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
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

	_, err = tx.Exec(ctx, "DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2", followerId, followedId)
	if err != nil {
		return fmt.Errorf("failed to delete follower: %w", err)
	}

	return nil
}

func (r *followerRepositoryImpl) userFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
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
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE follower_id = $1 AND followed_id = $2)", followerId, followedId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user follows user: %w", err)
	}

	return exists, nil
}

func (r *followerRepositoryImpl) getFollowers(ctx context.Context, userId uuid.UUID) ([]shared.User, error) {
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

	rows, err := tx.Query(ctx, "SELECT u.id, u.username, u.full_name, u.avatar FROM followers f JOIN users u ON f.follower_id = u.id WHERE f.followed_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to select followers: %w", err)
	}
	defer rows.Close()

	var userFollowers []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userFollowers = append(userFollowers, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return userFollowers, nil
}

func (r *followerRepositoryImpl) getFollowed(ctx context.Context, userId uuid.UUID) ([]shared.User, error) {
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

	rows, err := tx.Query(ctx, "SELECT u.id, u.username, u.full_name, u.avatar FROM followers f JOIN users u ON f.followed_id = u.id WHERE f.follower_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to select followed: %w", err)
	}
	defer rows.Close()

	var userFollowed []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userFollowed = append(userFollowed, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return userFollowed, nil
}
