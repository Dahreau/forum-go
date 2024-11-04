package database

import (
	"forum-go/internal/models"
	"math"
	"math/rand"
	"strconv"
)

func (s *service) GetComment() ([]models.Comment, error) {
	rows, err := s.db.Query("SELECT * FROM Comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.CommentId)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s *service) AddComment(name string) error {
	comment := models.Comment{
		CommentId: strconv.Itoa(rand.Intn(math.MaxInt32)),
		Name:       name,
	}
	query := "INSERT INTO Comment (comment_id,name) VALUES (?,?)"
	_, err := s.db.Exec(query, comment.CommentId, comment.Name)
	return err
}

func (s *service) DeleteComment(id string) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Delete comment post by id
	query := "DELETE FROM Comment WHERE comment_id=?"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) EditCategory(id, name string) error {
	query := "UPDATE Category SET name=? WHERE category_id=?"
	_, err := s.db.Exec(query, name, id)
	return err
}
