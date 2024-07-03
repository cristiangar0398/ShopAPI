package repository

import (
	"context"

	"github.com/cristiangar0398/ShopAPI/models"
)

var (
	implementation Repository
)

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctc context.Context, post *models.Post) error
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, id string, userID string) error
	ListPost(ctx context.Context, list uint64) ([]*models.Post, error)
	Close() error
}

func SetRepository(repository Repository) {
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementation.GetPostById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func UpdatePost(ctx context.Context, post *models.Post) error {
	return implementation.UpdatePost(ctx, post)
}

func DeletePost(ctx context.Context, id string, userID string) error {
	return implementation.DeletePost(ctx, id, userID)
}

func ListPost(ctx context.Context, list uint64) ([]*models.Post, error) {
	return implementation.ListPost(ctx, list)
}

func Close() error {
	return implementation.Close()
}
