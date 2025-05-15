package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/config"
	authMiddleware "github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/middleware"
	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/store"
)

// Server holds the HTTP server and all dependencies
type Server struct {
	cfg        *config.Config
	handler    *UserHandler
	userStore  store.UserStore
	authMid    *authMiddleware.AuthMiddleware
	httpServer *http.Server
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

	// Initialize auth middleware
	authMid := authMiddleware.NewAuthMiddleware(cfg.JWTSecretKey)

	// Initialize handler
	handler := NewUserHandler(userStore)

	// Create server
	server := &Server{
		cfg:        cfg,
		handler:    handler,
		userStore:  userStore,
		authMid:    authMid,
		httpServer: nil,
	}

	return server, nil
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
		AllowedOrigins:   []string{"*"}, // In production, you should specify actual origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum cache age for preflight options request
	}))

	// Public routes - no authentication needed
	// None for this service

	// Protected routes - authentication required
	r.Group(func(r chi.Router) {
		// Apply authentication middleware to all routes in this group
		r.Use(s.authMid.Authenticate)

		// User profile routes
		r.Route("/users", func(r chi.Router) {
			// Get user by ID - requires the user to be the owner or an admin
			r.With(s.authMid.RequireOwnerOrAdmin).Get("/{id}", s.handler.GetUserHandler)

			// Update user - requires the user to be the owner or an admin
			r.With(s.authMid.RequireOwnerOrAdmin).Put("/{id}", s.handler.UpdateUserHandler)

			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(s.authMid.RequireAdmin)

				// List all users - admin only
				r.Get("/", s.handler.ListUsersHandler)

				// Update user role - admin only
				r.Put("/{id}/role", s.handler.UpdateRoleHandler)
			})
		})
	})

	return r
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
