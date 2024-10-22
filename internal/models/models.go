package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserId         string
	Username       string
	Email          string
	Password       string
	Role           string
	CreationDate   time.Time
	SessionId      sql.NullString
	SessionExpired sql.NullTime
}
