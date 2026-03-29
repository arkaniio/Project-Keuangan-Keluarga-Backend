package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"project-keuangan-keluarga/controller"
	"project-keuangan-keluarga/middleware"
)

// Setup creates the chi router, registers middleware, and mounts all routes.
func UserRoutes(userCtrl *controller.ControllerHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger) // custom structured JSON logger
	r.Use(chimw.Recoverer)   // recover from panics
	r.Use(chimw.RequestID)   // inject X-Request-Id header

	// Health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", userCtrl.Register)
			r.Post("/login", userCtrl.Login)
		})
	})

	return r
}
