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
	category := models.Category{
		CategoryId: strconv.Itoa(rand.Intn(math.MaxInt32)),
		Name:       name,
	}
	query := "INSERT INTO Category (category_id,name) VALUES (?,?)"
	_, err := s.db.Exec(query, category.CategoryId, category.Name)
	return err
}

func (s *service) DeleteCategory(id string) error {
	query := "DELETE FROM Category WHERE category_id=?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *service) EditCategory(id, name string) error {
	query := "UPDATE Category SET name=? WHERE category_id=?"
	_, err := s.db.Exec(query, name, id)
	return err
}
