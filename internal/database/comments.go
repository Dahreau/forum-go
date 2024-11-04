package database

import (
	"forum-go/internal/models"
	"math"
	"math/rand"
	"strconv"
)

func (s *service) GetComments(post models.Post) ([]models.Comment, error) {
	rows, err := s.db.Query(`
		SELECT c.comment_id, c.content, c.creation_date, c.user_id, c.post_id
		FROM Comment c
		WHERE c.post_id = ?`, post.PostId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.CommentId, &comment.Content, &comment.CreationDate, &comment.UserID, &comment.PostID)
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

func (s *service) EditComment(id, content string) error {
	query := "UPDATE Comment SET content=? WHERE comment_id=?"
	_, err := s.db.Exec(query, content, id)
	return err
}
