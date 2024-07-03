package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/cristiangar0398/ShopAPI/models"
)

type PostgresRepository struct {
	db *sql.DB
}

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

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id , post_content , user_id) VALUES ($1, $2 ,$3)", post.Id, post.Post_content, post.UserId)
	return err
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 and user_id = $3", post.Post_content, post.Id, post.UserId)
	return err
}
func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
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

func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	rows, err := repo.db.QueryContext(ctx, "SELECT id , post_content , created_at , user_id FROM posts WHERE id = $1", id)

	if err != nil {
		log.Fatal(err)
		return &post, nil
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.Post_content, &post.Created_at, &post.UserId); err != nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
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

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userdID string) error {
	fmt.Println(userdID)
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 and user_id = $2", id, userdID)
	return err
}

func (repo *PostgresRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts LIMIT $1 OFFSET $2", 5, page*5)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var posts []*models.Post
	for rows.Next() {
		var post = models.Post{}
		if err = rows.Scan(&post.Id, &post.Post_content, &post.UserId, &post.Created_at); err == nil {
			posts = append(posts, &post)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
