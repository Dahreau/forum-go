package database

import (
	"context"
	"database/sql"
	"fmt"
	"forum-go/internal/models"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	CreateUser(user models.User) error
	GetUsers() ([]models.User, error)
	GetUser(email, password string) (models.User, error)
	Close() error

	FindEmailUser(email string) (bool, error)
	FindUsername(username string) (bool, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) error

	FindUserCookie(cookie string) (models.User, error)

	GetPosts() ([]models.Post, error)
	GetPost(id string) (models.Post, error)
	AddPost(post models.Post, categories []models.Category) error
	DeletePost(id string) error
	EditPost(id, title string) error

	GetCategories() ([]models.Category, error)
	AddCategory(name string) error
	DeleteCategory(id string) error
	EditCategory(id, name string) error
}

type service struct {
	db *sql.DB
}

var (
	dburl      = "./db.sqlite"
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	// Read the SQL file
	file, err := os.Open("query.sql")
	if err != nil {
		log.Fatal("Error opening SQL file:", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}
	defer file.Close()
	// Read the contents of the SQL file
	sqlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading SQL file:", err)
	}
	sqlStatements := string(sqlBytes)

	// Execute the SQL statements in the file
	_, err = db.Exec(sqlStatements)
	if err != nil {
		log.Fatal("Error executing SQL script:", err)

	}

	fmt.Println("Database initialized successfully!")

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}
