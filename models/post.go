package models

import "time"

type Products struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json"image_url"`
	Price       int       `json:"price"`
	Created_at  time.Time `json:"created_at"`
	UserId      string    `json:"userId"`
}
