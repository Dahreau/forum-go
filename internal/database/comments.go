package database

import (
	"forum-go/internal/models"
)

func (s *service) GetComments(post models.Post) ([]models.Comment, error) {
	// Query to get all comments for a post
	rows, err := s.db.Query(`
        SELECT c.comment_id, c.content, c.creation_date, c.user_id, c.post_id, u.username
        FROM Comment c
        JOIN User u ON c.user_id = u.user_id
        WHERE c.post_id = ?
        ORDER BY c.creation_date ASC`, post.PostId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.CommentId, &comment.Content, &comment.CreationDate, &comment.UserID, &comment.PostID, &comment.Username)
		if err != nil {
			return nil, err
		}
		// Format creation date
		comment.FormattedCreationDate = comment.CreationDate.Format("02/01/06 - 15:04")
		comments = append(comments, comment)
	}
	// Get user likes for each comment
	for i := range comments {
		userlikes, err := s.GetCommentLikes(comments[i].CommentId)
		if err != nil {
			return nil, err
		}
		comments[i].UserLikes = userlikes
		comments[i].Likes, comments[i].Dislikes = s.GetLikesCount(userlikes)
	}
	return comments, nil
}

func (s *service) AddComment(comment models.Comment) error {
	// Query insert all fields in comment table
	query := "INSERT INTO Comment (comment_id,content, creation_date, user_id, post_id) VALUES (?,?,?,?,?)"
	_, err := s.db.Exec(query, comment.CommentId, comment.Content, comment.CreationDate, comment.UserID, comment.PostID)
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
	// Query to update comment content
	query := "UPDATE Comment SET content=? WHERE comment_id=?"
	_, err := s.db.Exec(query, content, id)
	return err
}
