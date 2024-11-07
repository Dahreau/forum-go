package database

import (
	"forum-go/internal/models"
	"strings"
)

func (s *service) GetPosts() ([]models.Post, error) {
	// Query to retrieve posts with concatenated category IDs and names
	rows, err := s.db.Query(`
		SELECT 
			p.post_id, 
			p.title, 
			p.content, 
			p.user_id, 
			p.creation_date, 
			p.update_date, 
			GROUP_CONCAT(c.category_id) AS category_ids, 
			GROUP_CONCAT(c.name) AS category_names 
		FROM 
			Post p
		LEFT JOIN 
			Post_Category pc ON p.post_id = pc.post_id
		LEFT JOIN 
			Category c ON pc.category_id = c.category_id
		GROUP BY 
			p.post_id
		ORDER BY 
			p.creation_date DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve all users
	users, err := s.GetUsers()
	if err != nil {
		return nil, err
	}

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var categoryIDs, categoryNames string
		err := rows.Scan(&post.PostId, &post.Title, &post.Content, &post.UserID, &post.CreationDate, &post.UpdateDate, &categoryIDs, &categoryNames)
		if err != nil {
			return nil, err
		}
		post.FormattedCreationDate = post.CreationDate.Format("Jan 02, 2006 - 15:04:05")

		// Parse concatenated category IDs and names
		categoryIdList := strings.Split(categoryIDs, ",")
		categoryNameList := strings.Split(categoryNames, ",")

		for i := range categoryIdList {
			category := models.Category{
				CategoryId: categoryIdList[i],
				Name:       categoryNameList[i],
			}
			post.Categories = append(post.Categories, category)
		}

		// Attach user information
		for _, user := range users {
			if user.UserId == post.UserID {
				post.User = user
				break
			}
		}

		// Fetch comments and likes
		post.Comments, err = s.GetComments(post)
		if err != nil {
			return nil, err
		}
		post.NbOfComments = len(post.Comments)

		userLikes, err := s.GetPostLikes(post.PostId)
		if err != nil {
			return nil, err
		}
		post.UserLikes = userLikes
		post.Likes, post.Dislikes = s.GetLikesCount(userLikes)

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *service) GetPost(id string) (models.Post, error) {
	post := models.Post{}
	user := models.User{}
	query := `
		SELECT p.post_id, p.title, p.content, p.user_id, p.creation_date, p.update_date, 
			   u.user_id, u.username, u.email,
			   c.category_id, c.name
		FROM Post p 
		JOIN User u ON p.user_id = u.user_id 
		LEFT JOIN Post_Category pc ON p.post_id = pc.post_id
		LEFT JOIN Category c ON pc.category_id = c.category_id
		WHERE p.post_id = ?`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return post, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&post.PostId, &post.Title, &post.Content, &post.UserID, &post.CreationDate, &post.UpdateDate,
			&user.UserId, &user.Username, &user.Email,
			&category.CategoryId, &category.Name,
		)
		if err != nil {
			return post, err
		}
		categories = append(categories, category)
	}
	post.Comments, err = s.GetComments(post)
	if err != nil {
		return post, err
	}
	userlikes, err := s.GetPostLikes(post.PostId)
	if err != nil {
		return post, err
	}
	post.UserLikes = userlikes
	post.Likes, post.Dislikes = s.GetLikesCount(userlikes)
	post.FormattedCreationDate = post.CreationDate.Format("Jan 02, 2006 - 15:04:05")
	post.User = user
	post.Categories = categories
	return post, nil
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
