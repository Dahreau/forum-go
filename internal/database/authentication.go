package database

import (
	"forum-go/internal/models"
)

func (s *service) CreateUser(User models.User) error {

	query :="INSERT INTO User (user_id, email, username, password, role, creation_date, session_id, session_expire) VALUES (?, ?, ?, ?, ?, ?,?,?)"
	_, err := s.db.Exec(query, User.UserId, User.Email, User.Username, User.Password, User.Role, User.CreationDate, User.SessionId, User.SessionExpired)
	return err
}
