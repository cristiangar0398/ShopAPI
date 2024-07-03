package models

import "time"

type Post struct {
	Id           string    `json:"id"`
	Post_content string    `json:"post_content"`
	Created_at   time.Time `json:"created_at"`
	UserId       string    `json:"userId"`
}
