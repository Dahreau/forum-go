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
		post.FormattedCreationDate = post.CreationDate.Format("Jan 02, 2006 - 15:04:05")
		posts = append(posts, post)
	}
	return posts, nil
}
func (s *service) GetPost(id string) (models.Post, error) {
	post := models.Post{}
	user := models.User{}
	query := `
		SELECT p.post_id, p.title, p.content, p.user_id, p.creation_date, p.update_date, 
			   u.user_id, u.username, u.email 
		FROM Post p 
		JOIN User u ON p.user_id = u.user_id 
		WHERE p.post_id = ?`
	err := s.db.QueryRow(query, id).Scan(
		&post.PostId, &post.Title, &post.Content, &post.UserID, &post.CreationDate, &post.UpdateDate,
		&user.UserId, &user.Username, &user.Email,
	)
	post.FormattedCreationDate = post.CreationDate.Format("Jan 02, 2006 - 15:04:05")
	post.User = user
	return post, err
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
