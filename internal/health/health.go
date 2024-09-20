package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/uptrace/bun"
)

type healthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type health struct {
	db  *bun.DB
	ctx context.Context
}

func (h *health) healthHandler(w http.ResponseWriter, r *http.Request) {
	err := h.db.PingContext(h.ctx)

	res := healthResponse{
		Status:  "Up",
		Message: "All systems operational!",
	}

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		res.Status = "Down"
		res.Message = "Unable to connect to database"

		jsonRes, err := json.Marshal(res)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonRes)
		return
	}

	jsonRes, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func HealthRouteGroup(mux *http.ServeMux, ctx context.Context, db *bun.DB) {
	healthService := &health{db, ctx}
	mux.HandleFunc("GET /health", healthService.healthHandler)
}
