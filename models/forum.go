package models

import "time"

type Forum struct {
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"user"`
	Slug    string `json:"slug" db:"slug"`
	Posts   int    `json:"posts" db:"posts"`
	Threads int    `json:"threads" db:"threads"`
}

type User struct {
	Nickname string `json:"nickname,omitempty" db:"nickname"`
	Fullname string `json:"fullname,omitempty" db:"fullname"`
	About    string `json:"about,omitempty" db:"about"`
	Email    string `json:"email,omitempty" db:"email"`
}

type MessageStatus struct {
	Message string `json:"message,omitempty" db:"message"`
}

type Thread struct {
	Id      int       `json:"id,omitempty" db:"id"`
	Title   string    `json:"title" db:"title"`
	Author  string    `json:"author" db:"author"`
	Forum   string    `json:"forum" db:"forum"`
	Message string    `json:"message" db:"message"`
	Votes   int       `json:"votes" db:"votes"`
	Slug    string    `json:"slug" db:"slug"`
	Created time.Time `json:"created" db:"created"`
}

type Post struct {
	Id       int       `json:"id" db:"id"`
	Parent   int       `json:"parent" db:"parent"`
	Author   string    `json:"author" db:"author"`
	Message  string    `json:"message" db:"message"`
	IsEdited bool      `json:"isEdited,omitempty" db:"isedited"`
	Forum    string    `json:"forum" db:"forum"`
	Thread   int       `json:"thread" db:"thread"`
	Created  time.Time `json:"created" db:"created"`
}