package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserId        string         `db:"user_id"`
	Email         string         `db:"email"`
	Username      string         `db:"username"`
	Password      string         `db:"password"`
	Role          string         `db:"role"`
	CreationDate  time.Time      `db:"creation_date"`
	SessionId     sql.NullString `db:"session_id"`
	SessionExpire sql.NullTime   `db:"session_expire"`
	Posts         []Post         `db:"-"`
}

type Category struct {
	CategoryId string `db:"category_id"`
	Name       string `db:"name"`
}

type Post struct {
	PostId                string       `db:"post_id"`
	Title                 string       `db:"title"`
	Content               string       `db:"content"`
	UserID                string       `db:"user_id"`
	CreationDate          time.Time    `db:"creation_date"`
	FormattedCreationDate string       `db:"-"`
	UpdateDate            sql.NullTime `db:"update_date"`
	User                  User         `db:"-"`
	Categories            []Category   `db:"-"`
	Comments              []Comment    `db:"-"`
	NbOfComments          int          `db:"-"`
	Likes                 int          `db:"-"`
	Dislikes              int          `db:"-"`
}

type Comment struct {
	CommentId             string    `db:"comment_id"`
	Content               string    `db:"content"`
	CreationDate          time.Time `db:"creation_date"`
	FormattedCreationDate string    `db:"-"`
	UserID                string    `db:"user_id"`
	PostID                string    `db:"post_id"`
	Username              string    `db:"-"`
	Likes                 int       `db:"-"`
	Dislikes              int       `db:"-"`
}

type PostCategory struct {
	PostId     string `db:"post_id"`
	CategoryId string `db:"category_id"`
}

type UserLike struct {
	LikeId    string `db:"like_id"`
	UserId    string `db:"user_id"`
	PostId    string `db:"post_id"`
	CommentId string `db:"comment_id"`
	IsLike    bool   `db:"is_like"`
}
