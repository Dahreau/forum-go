package server

import (
	"crypto/tls"
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
	posts      []models.Post
	SESSION_ID string
}

func NewServer() *http.Server {

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("error cerst", err)
	}

	// Configurer les options TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Force TLS v1.2 ou supérieur
		// Configurer les suites de chiffrement (optionnel)
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true, // Prioriser les suites de chiffrement du serveur
	}

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
	posts, err := NewServer.db.GetPosts()
	if err != nil {
		fmt.Println("Error getting posts: ", err)
	} else {
		NewServer.posts = posts
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		TLSConfig:    tlsConfig,
	}

	return server
}
