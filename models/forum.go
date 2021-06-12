package models

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