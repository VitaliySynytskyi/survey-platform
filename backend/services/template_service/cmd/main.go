package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Application name for logging
	appName := "template-service"
	log.Printf("Starting %s...", appName)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", defaultHandler)

	// Configure the HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("%s is listening on port %s", appName, port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt/terminate signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-shutdown:
		log.Printf("Shutdown signal received, shutting down gracefully...")

		// Create a deadline for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Attempt to gracefully shutdown the server
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Error during server shutdown: %v", err)
			err = server.Close()
		}

		// If we get here, the server is shutting down gracefully
		log.Printf("Server gracefully stopped")
	}
}

// healthHandler responds to health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}

// defaultHandler responds to all other requests
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Return 404 if not the root path
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"Not Found"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Template Service API"}`)
}
