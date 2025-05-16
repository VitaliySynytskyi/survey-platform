package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/auth"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/config"
	authMiddleware "github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/middleware"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/store"
)

// Server holds the HTTP server and all dependencies
type Server struct {
	cfg          *config.Config
	handler      *Handler
	userStore    store.UserStore
	tokenManager *auth.TokenManager
	httpServer   *http.Server
}

// HealthStatus represents the status of the service and its dependencies
type HealthStatus struct {
	Status      string `json:"status"`
	Database    string `json:"database,omitempty"`
	Details     string `json:"details,omitempty"`
	ServiceName string `json:"service_name"`
}

// NewServer creates a new server with all dependencies
func NewServer(cfg *config.Config) (*Server, error) {
	// Create connection string for PostgreSQL
	connString := "postgres://" + cfg.Database.User + ":" + cfg.Database.Password +
		"@" + cfg.Database.Host + ":" + cfg.Database.Port +
		"/" + cfg.Database.DBName + "?sslmode=" + cfg.Database.SSLMode

	// Initialize store
	userStore, err := store.NewPostgresStore(connString)
	if err != nil {
		return nil, err
	}

	// Initialize JWT token manager
	tokenManager := auth.NewTokenManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExp,
		cfg.JWT.RefreshTokenExp,
	)

	// Initialize handler
	handler := NewHandler(userStore, tokenManager)

	// Create server
	server := &Server{
		cfg:          cfg,
		handler:      handler,
		userStore:    userStore,
		tokenManager: tokenManager,
		httpServer:   nil,
	}

	return server, nil
}

// SetupDatabase ensures the database schema is ready
func (s *Server) SetupDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.userStore.EnsureSchema(ctx); err != nil {
		return err
	}

	return nil
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:80", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum cache age for preflight options request
	}))

	// Health check endpoint
	r.Get("/health", s.healthCheckHandler)

	// Auth routes - no authentication needed
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.handler.RegisterHandler)
		r.Post("/login", s.handler.LoginHandler)
		r.Post("/refresh", s.handler.RefreshHandler)

		// Protected route - requires authentication
		authMid := authMiddleware.NewAuthMiddleware(s.tokenManager)
		r.With(authMid.Authenticate).Get("/me", s.handler.MeHandler)
	})

	return r
}

// healthCheckHandler checks if the service and its dependencies are healthy
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := HealthStatus{
		Status:      "healthy",
		ServiceName: "auth_service",
	}

	// Check database connection
	if err := s.userStore.Ping(ctx); err != nil {
		log.Printf("Health check failed: database connection error: %v", err)
		status.Status = "unhealthy"
		status.Database = "disconnected"
		status.Details = err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status.Database = "connected"
		w.WriteHeader(http.StatusOK)
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Set up all routes
	router := s.SetupRoutes()

	// Configure the HTTP server
	s.httpServer = &http.Server{
		Addr:         ":" + s.cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", s.cfg.Server.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Close the database connection
	s.userStore.Close()

	log.Println("Server stopped gracefully")
}

// WaitForSignal blocks until an OS interrupt signal is received
func (s *Server) WaitForSignal() {
	// Create a channel to listen for OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-c
	log.Println("Received shutdown signal")
}
