package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/cristiangar0398/ShopAPI/models"
)

type PostgresRepository struct {
	db *sql.DB
}

var (
	product models.Products
	user    models.User
)

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id , email , password) VALUES ($1, $2 ,$3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) InsertProduct(ctx context.Context, product *models.Products) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO products (id , title , description , image_url , price , user_id) VALUES ($1, $2 ,$3 ,$4 ,$5 ,$6)", product.Id, product.Title, product.Description, product.ImageUrl, product.Price, product.UserId)
	return err
}

func (repo *PostgresRepository) UpdateProduct(ctx context.Context, product *models.Products) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE products SET title = $1, description = $2, image_url = $3, price = $4 WHERE id = $5 and user_id = $6", product.Title, product.Description, product.ImageUrl, product.Price, product.Id, product.UserId)
	return err
}
func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id , email FROM users WHERE id = $1", id)

	if err != nil {
		log.Fatal(err)
		return &user, nil
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err != nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) GetProductById(ctx context.Context, id string) (*models.Products, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id , title , description , image_url , price ,created_at,user_id FROM products WHERE id = $1", id)

	if err != nil {
		log.Fatal(err)
		return &product, nil
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&product.Id, &product.Title, &product.Description, &product.ImageUrl, &product.Price, &product.Created_at, &product.UserId); err != nil {
			return &product, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email , password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Fatal(cerr)
		}
	}()

	if rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		return &user, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (repo *PostgresRepository) DeleteProduct(ctx context.Context, id string, userdID string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1 and user_id = $2", id, userdID)
	return err
}

func (repo *PostgresRepository) ListProducts(ctx context.Context, page uint64) ([]*models.Products, error) {

	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM products LIMIT $1 OFFSET $2", 5, page*5)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var products []*models.Products
	for rows.Next() {
		if err = rows.Scan(&product.Id, &product.Title, &product.Description, &product.ImageUrl, &product.Price, &product.Created_at, &product.UserId); err == nil {
			products = append(products, &product)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
