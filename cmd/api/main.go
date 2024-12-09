package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"forum-go/internal/server"
	"forum-go/security"
)

func gracefulShutdown(apiServer *http.Server) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
}

func main() {
	go security.CleanupInactiveClients()
	//dotenv.Define(".env.prod")
	//shared.GoogleRedirectURL = dotenv.GetEnv("googleRedirectURL")
	//shared.GoogleClientSecret = dotenv.GetEnv("googleClientSecret")
	//shared.GoogleClientID = dotenv.GetEnv("googleClientID")
	//
	server := server.NewServer()

	go gracefulShutdown(server)

	fmt.Println("Server started on port", server.Addr)
	fmt.Println("https://localhost:8080")
	// err := server.ListenAndServe()
	err := server.ListenAndServeTLS("./cert.pem", "./key.pem")
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
