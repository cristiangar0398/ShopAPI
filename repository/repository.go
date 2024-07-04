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
	InsertProduct(ctc context.Context, product *models.Products) error
	GetProductById(ctx context.Context, id string) (*models.Products, error)
	UpdateProduct(ctx context.Context, product *models.Products) error
	DeleteProduct(ctx context.Context, id string, userID string) error
	ListProducts(ctx context.Context, list uint64) ([]*models.Products, error)
	Close() error
}

func SetRepository(repository Repository) {
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func InsertProduct(ctx context.Context, post *models.Products) error {
	return implementation.InsertProduct(ctx, post)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

func GetProductById(ctx context.Context, id string) (*models.Products, error) {
	return implementation.GetProductById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func UpdateProduct(ctx context.Context, post *models.Products) error {
	return implementation.UpdateProduct(ctx, post)
}

func DeleteProduct(ctx context.Context, id string, userID string) error {
	return implementation.DeleteProduct(ctx, id, userID)
}

func ListProducts(ctx context.Context, list uint64) ([]*models.Products, error) {
	return implementation.ListProducts(ctx, list)
}

func Close() error {
	return implementation.Close()
}
