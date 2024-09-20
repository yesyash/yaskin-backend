package main

import (
	"context"

	"github.com/yesyash/yaskin-backend/internal/config"
	"github.com/yesyash/yaskin-backend/internal/database"
	"github.com/yesyash/yaskin-backend/internal/logger"
	"github.com/yesyash/yaskin-backend/internal/server"
)

func main() {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a database connection
	db := database.New()
	defer db.Close()

	// Create a new server
	httpServer := server.NewServer(ctx, db)
	logger.Info("Server running on port", config.Port)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
