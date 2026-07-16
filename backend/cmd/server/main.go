package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tablelink-backend/internal/config"
	"tablelink-backend/internal/server"
)

func main() {
	// ---------------------------------------------------------------
	// Load configuration from environment / .env
	// ---------------------------------------------------------------
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ---------------------------------------------------------------
	// Build and start the server (DI + DB + routes all wired inside)
	// ---------------------------------------------------------------
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// ---------------------------------------------------------------
	// Graceful shutdown on SIGINT / SIGTERM
	// ---------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")
	if err := srv.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
