package server

import (
	"fmt"
	"net/http"
	"time"

	"forum-go/internal/database"
	"forum-go/internal/models"
)

type Server struct {
	port       int
	db         database.Service
	users      []models.User
	categories []models.Category
	SESSION_ID string
}

func NewServer() *http.Server {
	NewServer := &Server{
		port:       8080,
		db:         database.New(),
		SESSION_ID: "sRpyIJS9Zmerlpcpqhc1B0xxG7w6Gk1b",
	}
	users, err := NewServer.db.GetUsers()
	if err != nil {
		fmt.Println("Error getting users: ", err)
	} else {
		NewServer.users = users
	}
	categories, err := NewServer.db.GetCategories()
	if err != nil {
		fmt.Println("Error getting categories: ", err)
	} else {
		NewServer.categories = categories
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
