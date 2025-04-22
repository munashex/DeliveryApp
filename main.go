package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/munashex/goweb/config"
)

func main() {
	// ========================================
	// Initialize Configuration
	// ========================================
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log loaded configuration (remove in production if sensitive)
	log.Printf("Starting service with config: %+v", cfg.Sanitized())

	// ========================================
	// Database Connection (now using config)
	// ========================================
	// db, err := database.Connect(cfg.Database)
	// if err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer db.Close()

	// ========================================
	// Initialize HTTP Server with Config
	// ========================================
	router := http.NewServeMux()
	
	// Register routes
	// routes.RegisterRoutes(router, cfg)
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// ========================================
	// Start HTTP Server with Configurable Options
	// ========================================
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		switch cfg.Server.Protocol {
		case "http":
			serverErrors <- server.ListenAndServe()
		case "https":
			serverErrors <- server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile)
		default:
			serverErrors <- fmt.Errorf("invalid server protocol: %s", cfg.Server.Protocol)
		}
	}()

	// ========================================
	// Shutdown Logic with Configurable Timeout
	// ========================================
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case <-shutdown:
		log.Println("Starting graceful shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Fatalf("Could not stop server: %v", err)
			}
		}
	}

	log.Println("Server stopped")
}