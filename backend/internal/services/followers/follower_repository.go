package followers

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"
	"y_net/internal/services/users"

	"github.com/google/uuid"
)

type followerRepositoryI interface {
	follow(followerId uuid.UUID, followedId uuid.UUID) error
	unfollow(followerId uuid.UUID, followedId uuid.UUID) error
	getFollowers(userId uuid.UUID) ([]users.User, error)
	getFollowed(userId uuid.UUID) ([]users.User, error)
}

type followerRepository struct{}

func (i *followerRepository) follow(followerId uuid.UUID, followedId uuid.UUID) error {
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

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2)",
		followerId, followedId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert follower: %w", err)
	}

	return nil
}

func (i *followerRepository) unfollow(followerId uuid.UUID, followedId uuid.UUID) error {
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

	_, err = tx.Exec(context.Background(), "DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2", followerId, followedId)
	if err != nil {
		return fmt.Errorf("failed to delete follower: %w", err)
	}

	return nil
}

func (i *followerRepository) getFollowers(userId uuid.UUID) ([]users.User, error) {
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

	rows, err := tx.Query(context.Background(), "SELECT u.id, u.username, u.full_name, u.avatar FROM followers f JOIN users u ON f.follower_id = u.id WHERE f.followed_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to select followers: %w", err)
	}
	defer rows.Close()

	var userFollowers []users.User
	for rows.Next() {
		var user users.User
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

func (i *followerRepository) getFollowed(userId uuid.UUID) ([]users.User, error) {
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

	rows, err := tx.Query(context.Background(), "SELECT u.id, u.username, u.full_name, u.avatar FROM followers f JOIN users u ON f.followed_id = u.id WHERE f.follower_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to select followed: %w", err)
	}
	defer rows.Close()

	var userFollowed []users.User
	for rows.Next() {
		var user users.User
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
