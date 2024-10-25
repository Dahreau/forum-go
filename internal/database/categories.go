package database

import (
	"forum-go/internal/models"
	"math"
	"strconv"

	"math/rand"
)

func (s *service) GetCategories() ([]models.Category, error) {
	rows, err := s.db.Query("SELECT * FROM Category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.CategoryId, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (s *service) AddCategory(name string) error {
	query := "INSERT INTO Category (category_id,name) VALUES (?,?)"
	_, err := s.db.Exec(query, strconv.Itoa(rand.Intn(math.MaxInt32)), name)
	return err
}
