package database

import (
	"database/sql"
	"forum-go/internal/models"
	"forum-go/internal/shared"
)

func (s *service) Vote(postID, commentID, userID string, isLike bool) error {
	// Check if user has already liked or disliked the post
	var row *sql.Row
	var userlike models.UserLike
	if commentID != "''" && commentID != "" {
		query := "SELECT * FROM User_like WHERE comment_id=? AND user_id=?"
		row = s.db.QueryRow(query, commentID, userID)
	} else {
		query := "SELECT * FROM User_like WHERE post_id=? AND user_id=? AND comment_id = ''"
		row = s.db.QueryRow(query, postID, userID)
	}
	if err := row.Scan(&userlike.LikeId, &userlike.IsLike, &userlike.UserId, &userlike.PostId, &userlike.CommentId); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		userlike.LikeId = ""
	}

	if userlike.LikeId == "" {
		// If user has not liked or disliked the post, insert the like
		userlike.LikeId = shared.ParseUUID(shared.GenerateUUID())
		query := "INSERT INTO User_like (like_id, user_id, post_id, comment_id, isLiked) VALUES (?,?,?,?,?)"
		_, err := s.db.Exec(query, userlike.LikeId, userID, postID, commentID, isLike)
		return err
	}

	if userlike.IsLike == isLike {
		// If user has already liked or disliked the post, delete the like
		_, err := s.db.Exec("DELETE FROM User_like WHERE like_id=?", userlike.LikeId)
		return err

	}
	// If user has already liked or disliked the post, update the like
	query := "UPDATE User_like SET isLiked=? WHERE like_id=?"
	_, err := s.db.Exec(query, isLike, userlike.LikeId)
	return err
}

func (s *service) GetPostLikes(postID string) ([]models.UserLike, error) {
	// Query to get all likes for a post
	query := "SELECT * FROM User_like WHERE post_id=? AND comment_id = ''"
	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userlikes []models.UserLike
	for rows.Next() {
		var userlike models.UserLike
		if err := rows.Scan(&userlike.LikeId, &userlike.IsLike, &userlike.UserId, &userlike.PostId, &userlike.CommentId); err != nil {
			return nil, err
		}
		userlikes = append(userlikes, userlike)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userlikes, nil
}
func (s *service) GetCommentLikes(commentID string) ([]models.UserLike, error) {
	// Query to get all likes for a comment
	query := "SELECT * FROM User_like WHERE comment_id=?"
	rows, err := s.db.Query(query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userlikes []models.UserLike
	for rows.Next() {
		var userlike models.UserLike
		if err := rows.Scan(&userlike.LikeId, &userlike.IsLike, &userlike.UserId, &userlike.PostId, &userlike.CommentId); err != nil {
			return nil, err
		}
		userlikes = append(userlikes, userlike)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userlikes, nil
}

func (s *service) GetLikesCount(userlikes []models.UserLike) (int, int) {
	// Get the number of likes and dislikes for a post
	var likes, dislikes int
	for _, like := range userlikes {
		if like.IsLike {
			likes++
		} else {
			dislikes++
		}
	}
	return likes, dislikes
}

func (s *service) DeleteLikes(postID string) error {
	// Delete all likes for a post
	query := "DELETE FROM User_like WHERE post_id=?"
	_, err := s.db.Exec(query, postID)
	return err
}
func (s *service) DeleteCommentLikes(commentID string) error {
	// Delete all likes for a comment
	query := "DELETE FROM User_like WHERE comment_id=?"
	_, err := s.db.Exec(query, commentID)
	return err
}
