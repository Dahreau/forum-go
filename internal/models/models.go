package models

import "time"

type User struct {
	UserId         string
	Username       string
	Email          string
	Password       string
	Role           string
	CreationDate   time.Time
	SessionId      string
	SessionExpired time.Time
}
