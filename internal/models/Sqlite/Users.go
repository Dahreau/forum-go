package sqlite

import (
	"forum-go/internal/database"
)

type UserModel struct {
	Db *database.Service
}

func (u *UserModel) Insert(Email, Username, Password string) error {
	return nil
}
