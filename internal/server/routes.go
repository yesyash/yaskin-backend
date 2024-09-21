package server

import (
	"context"
	"net/http"

	"github.com/yesyash/yaskin-backend/internal/documents"
	"github.com/yesyash/yaskin-backend/internal/health"
)

func (s *Server) RegisterRoutes(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	health.HealthRouteGroup(mux, ctx, s.db)
	documents.DocumentGroup(mux, ctx, s.db)

	mux.HandleFunc("GET /public/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		// return the file from public folder which is in the root of the project
		http.ServeFile(w, r, "public/"+filename)
	})

	return mux
}
