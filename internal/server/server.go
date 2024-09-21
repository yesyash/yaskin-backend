package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"github.com/yesyash/yaskin-backend/internal/config"
)

type Server struct {
	db *bun.DB
}

func NewServer(ctx context.Context, db *bun.DB) *http.Server {
	NewServer := &Server{db}

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"http://localhost:3000"}, // Adjust this to match your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      c.Handler(NewServer.RegisterRoutes(ctx)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
