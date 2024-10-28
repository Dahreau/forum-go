package database

import (
	"forum-go/internal/models"
)

func (s *service) GetPosts() ([]models.Post, error) {
	rows, err := s.db.Query("SELECT * FROM Post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostId, &post.Title, &post.Content, &post.UserID, &post.CreationDate, &post.UpdateDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *service) AddPost(post models.Post) error {
	query := "INSERT INTO Post (post_id,title,content,user_id,creation_date) VALUES (?,?,?,?,?)"
	_, err := s.db.Exec(query, post.PostId, post.Title, post.Content, post.UserID, post.CreationDate)
	return err
}

func (s *service) DeletePost(id string) error {
	query := "DELETE FROM Post WHERE post_id=?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *service) EditPost(id, title string) error {
	query := "UPDATE Post SET title=? WHERE post_id=?"
	_, err := s.db.Exec(query, title, id)
	return err
}
