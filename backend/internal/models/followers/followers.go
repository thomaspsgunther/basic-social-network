package followers

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"
	"y_net/internal/models/users"

	"github.com/google/uuid"
)

type Follower struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

func Follow(followerId uuid.UUID, followedId uuid.UUID) error {
	follower := Follower{FollowerID: followerId, FollowedID: followedId}

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
		follower.FollowerID, follower.FollowedID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert follower: %w", err)
	}

	return nil
}

func Unfollow(followerId uuid.UUID, followedId uuid.UUID) error {
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

func GetFollowers(userId uuid.UUID) ([]users.User, error) {
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

	var followers []users.User
	for rows.Next() {
		var follower users.User
		if err := rows.Scan(&follower.ID, &follower.Username, &follower.FullName, &follower.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		followers = append(followers, follower)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return followers, nil
}

func GetFollowed(userId uuid.UUID) ([]users.User, error) {
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

	var followed []users.User
	for rows.Next() {
		var followedUser users.User
		if err := rows.Scan(&followedUser.ID, &followedUser.Username, &followedUser.FullName, &followedUser.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		followed = append(followed, followedUser)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return followed, nil
}
