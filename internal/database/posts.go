package database

import (
	"forum-go/internal/models"
)

func (s *service) GetPosts() ([]models.Post, error) {
	rows, err := s.db.Query(`
		SELECT p.post_id, p.title, p.content, p.user_id, p.creation_date, p.update_date, 
			   c.category_id, c.name 
		FROM Post p 
		LEFT JOIN Post_Category pc ON p.post_id = pc.post_id 
		LEFT JOIN Category c ON pc.category_id = c.category_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := s.GetUsers()
	if err != nil {
		return nil, err
	}
	postMap := make(map[string]*models.Post)
	for rows.Next() {
		var post models.Post
		var category models.Category
		err := rows.Scan(&post.PostId, &post.Title, &post.Content, &post.UserID, &post.CreationDate, &post.UpdateDate, &category.CategoryId, &category.Name)
		if err != nil {
			return nil, err
		}
		post.FormattedCreationDate = post.CreationDate.Format("Jan 02, 2006 - 15:04:05")
		if existingPost, ok := postMap[post.PostId]; ok {
			existingPost.Categories = append(existingPost.Categories, category)
		} else {
			for _, user := range users {
				if user.UserId == post.UserID {
					post.User = user
					break
				}
			}
			post.Categories = append(post.Categories, category)
			postMap[post.PostId] = &post
		}
	}
	posts := make([]models.Post, 0, len(postMap))
	for _, post := range postMap {
		post.NbOfComments = len(post.Comments)
		posts = append(posts, *post)
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

func (s *service) AddPost(post models.Post, categories []models.Category) error {
	query := "INSERT INTO Post (post_id,title,content,user_id,creation_date) VALUES (?,?,?,?,?)"
	_, err := s.db.Exec(query, post.PostId, post.Title, post.Content, post.UserID, post.CreationDate)
	if err != nil {
		return err
	}
	for _, category := range categories {
		query = "INSERT INTO Post_Category (post_id,category_id) VALUES (?,?)"
		_, err = s.db.Exec(query, post.PostId, category.CategoryId)
		if err != nil {
			return err
		}
	}
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
