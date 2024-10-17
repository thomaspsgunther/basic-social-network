package users

import (
	"context"
	"fmt"
	"time"
	database "y-net/internal/database/postgres"
	"y-net/internal/services/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Create(ctx context.Context, user shared.User) (uuid.UUID, error)
	Get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error)
	GetBySearch(ctx context.Context, searchStr string) ([]shared.User, error)
	GetPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error)
	Update(ctx context.Context, user shared.User, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	Unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	UserFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, id uuid.UUID) ([]shared.User, error)
	GetFollowed(ctx context.Context, id uuid.UUID) ([]shared.User, error)
}

type userUsecaseImpl struct {
	usecase    IUserUsecase
	repository iUserRepository
}

func NewUserUsecase() IUserUsecase {
	return &userUsecaseImpl{
		usecase:    &userUsecaseImpl{},
		repository: &userRepositoryImpl{},
	}
}

func (u *userUsecaseImpl) Create(ctx context.Context, user shared.User) (uuid.UUID, error) {
	if user.Username == "" || user.Password == "" {
		return uuid.Nil, fmt.Errorf("username and password must not be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	id, err := u.repository.create(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (u *userUsecaseImpl) Get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error) {
	users, err := u.repository.get(ctx, idList)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUsecaseImpl) GetBySearch(ctx context.Context, searchStr string) ([]shared.User, error) {
	users, err := u.repository.getBySearch(ctx, searchStr)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUsecaseImpl) GetPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
	posts, err := u.repository.getPostsFromUser(ctx, userId, limit, lastCreatedAt, lastId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (u *userUsecaseImpl) Update(ctx context.Context, user shared.User, id uuid.UUID) error {
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

	err := u.repository.update(ctx, user, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repository.delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) Follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	err := u.repository.follow(ctx, followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) Unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	err := u.repository.unfollow(ctx, followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) UserFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
	follows, err := u.repository.userFollowsUser(ctx, followerId, followedId)
	if err != nil {
		return false, err
	}

	return follows, nil
}

func (u *userUsecaseImpl) GetFollowers(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	users, err := u.repository.getFollowers(ctx, id)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUsecaseImpl) GetFollowed(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	users, err := u.repository.getFollowed(ctx, id)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func Authenticate(ctx context.Context, user shared.User) (bool, error) {
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
		database.HandleTransaction(ctx, tx, err)
	}()

	var hashedPassword string
	err = tx.QueryRow(ctx, "SELECT password FROM users WHERE username = $1", user.Username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	return checkPasswordHash(user.Password, hashedPassword), nil
}

// GetUsernameByUserID checks if a user exists in database by given id
func GetUsernameByUserID(ctx context.Context, id uuid.UUID) (string, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var username string
	err = tx.QueryRow(ctx, "SELECT username FROM users WHERE id = $1", id).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// GetUserIdByUsername checks if a user exists in database by given username
func GetUserIdByUsername(ctx context.Context, username string) (uuid.UUID, error) {
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
		database.HandleTransaction(ctx, tx, err)
	}()

	var id uuid.UUID
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&id)
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
