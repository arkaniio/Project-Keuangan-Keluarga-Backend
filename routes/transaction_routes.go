package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"project-keuangan-keluarga/controller"
	"project-keuangan-keluarga/middleware"
)

func KeuanganRoutes(keuanganCtrl *controller.ControllerHandlerKeuangan) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)         // custom structured JSON logger
	r.Use(middleware.MiddlewareAuth) // use
	r.Use(chimw.Recoverer)           // recover from panics
	r.Use(chimw.RequestID)           // inject X-Request-Id header

	// Health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// API v1 routes

	return r
}
