package models

import "time"

type Forum struct {
	Slug    string `json:"slug" db:"slug"`
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"user"`
	Posts   int    `json:"posts,omitempty" db:"posts"`
	Threads int    `json:"threads,omitempty" db:"threads"`
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
	Title   string    `json:"title,omitempty" db:"title"`
	Author  string    `json:"author,omitempty" db:"author"`
	Forum   string    `json:"forum,omitempty" db:"forum"`
	Message string    `json:"message,omitempty" db:"message"`
	Votes   int       `json:"votes,omitempty" db:"votes"`
	Slug    string    `json:"slug,omitempty" db:"slug"`
	Created time.Time `json:"created,omitempty" db:"created"`
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

type AllInfo struct {
	Post    *Post   `json:"post" db:"post"`
	Thread  *Thread    `json:"thread,omitempty" db:"thread"`
	User  *User `json:"user" db:"user,omitempty"`
	Forum *Forum `json:"forum,omitempty" db:"forum"`
}

type Voice struct {
	Thread  int    `json:"thread,omitempty" db:"thread"`
	Voice  int    `json:"voice,omitempty" db:"voice"`
	Nickname  string `json:"nickname" db:"nickname,omitempty"`
}

type Service struct {
	User  int `json:"user" db:"user,omitempty"`
	Forum  int    `json:"forum,omitempty" db:"forum"`
	Thread  int    `json:"thread,omitempty" db:"thread"`
	Post int `json:"post" db:"post,omitempty"`
}