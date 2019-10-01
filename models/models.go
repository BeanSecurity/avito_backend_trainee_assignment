package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UsersID   []int     `json:"users"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat"`
	AuthorID  int       `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
