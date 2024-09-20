package server

import (
	"context"
	"net/http"

	"github.com/yesyash/yaskin-backend/internal/health"
)

func (s *Server) RegisterRoutes(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	health.HealthRouteGroup(mux, ctx, s.db)

	return mux
}
