package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"project-keuangan-keluarga/controller"
	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/middleware/ratelimiter"
)

// CategoryRoutes creates the chi router for category-related endpoints.
func GoalsRoutes(goalsCtrl controller.ControllerGoals, generalLimiter *ratelimiter.Limiter) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)                              // custom structured JSON logger
	r.Use(middleware.MiddlewareAuth)                      // auth required
	r.Use(chimw.Recoverer)                                // recover from panics
	r.Use(chimw.RequestID)                                // inject X-Request-Id header
	r.Use(middleware.RateLimitMiddleware(generalLimiter)) // general rate limit: 60 req/min

	// Health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// API v1 routesww
	r.Post("/", goalsCtrl.CreateNewGoals_Bp)
	r.Get("/", goalsCtrl.GetAllGoals_Bp)
	r.Delete("/delete", goalsCtrl.DeleteGoals_Bp)
	r.Put("/update", goalsCtrl.UpdateGoals_Bp)
	r.Get("/progress", goalsCtrl.TrackingProgressGoals_Bp)

	return r
}
