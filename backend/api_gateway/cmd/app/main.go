package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/middleware"
	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/proxy"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create proxy router
	router, err := proxy.NewProxyRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to create proxy router: %v", err)
	}

	// Create a new middleware chain
	middlewareChain := middleware.New(router)

	// Add general middleware
	middlewareChain = middleware.AddLogging(middlewareChain)
	middlewareChain = middleware.AddCORS(middlewareChain, cfg.CORS)

	// Add rate limiting if enabled
	if cfg.RateLimiting.Enabled {
		middlewareChain = middleware.AddRateLimiting(middlewareChain, cfg.RateLimiting)
	}

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      middlewareChain,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting API Gateway on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
