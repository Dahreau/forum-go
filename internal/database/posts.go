package database

import (
	"forum-go/internal/models"
	"math"
	"math/rand"
	"strconv"
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
		err := rows.Scan(&post.PostId, &post.Title)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *service) AddPost(title string) error {
	post := models.Post{
		PostId: strconv.Itoa(rand.Intn(math.MaxInt32)),
		Title:  title,
	}
	query := "INSERT INTO Post (post_id,title) VALUES (?,?)"
	_, err := s.db.Exec(query, post.PostId, post.Title)
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
