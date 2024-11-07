package database

import (
	"database/sql"
	"forum-go/internal/models"
	"math"
	"math/rand"
	"strconv"
)

func (s *service) Vote(postID, commentID, userID string, isLike bool) error {
	var row *sql.Row
	var userlike models.UserLike
	if commentID != "''" {
		query := "SELECT * FROM User_like WHERE comment_id=? AND user_id=?"
		row = s.db.QueryRow(query, commentID, userID)
	} else {
		query := "SELECT * FROM User_like WHERE post_id=? AND user_id=?"
		row = s.db.QueryRow(query, postID, userID)
	}
	if err := row.Scan(&userlike.LikeId, &userlike.IsLike, &userlike.UserId, &userlike.PostId, &userlike.CommentId); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		userlike.LikeId = ""
	}

	if userlike.LikeId == "" {
		userlike.LikeId = strconv.Itoa(rand.Intn(math.MaxInt32))
		query := "INSERT INTO User_like (like_id, user_id, post_id, comment_id, isLiked) VALUES (?,?,?,?,?)"
		_, err := s.db.Exec(query, userlike.LikeId, userID, postID, commentID, isLike)
		return err
	}

	if userlike.IsLike == isLike {
		_, err := s.db.Exec("DELETE FROM User_like WHERE like_id=?", userlike.LikeId)
		return err

	}
	query := "UPDATE User_like SET isLiked=? WHERE like_id=?"
	_, err := s.db.Exec(query, isLike, userlike.LikeId)
	return err
}

func (s *service) GetPostLikes(postID string) ([]models.UserLike, error) {
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
