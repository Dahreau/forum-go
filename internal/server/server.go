package server

import (
	"fmt"
	"net/http"
	"time"

	"forum-go/internal/database"
)

type Server struct {
	port int
	db   database.Service
	//Users *sqlite.UserModel
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: 8080,
		db:   database.New(),
		// Users: &sqlite.UserModel{
		// 	Db: &db,
		// },
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