package routes

import (
	"encoding/json"
	"net/http"
	"project-keuangan-keluarga/controller"
	"project-keuangan-keluarga/middleware"
	"project-keuangan-keluarga/middleware/ratelimiter"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func FamilieRoutes(familieCtrl controller.ControllerHandlerFamilie, generalLimiter *ratelimiter.Limiter) *chi.Mux {
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

	// API v1 routes
	r.Post("/", familieCtrl.CreateNewFamilie_Bp)
	r.Delete("/:id", familieCtrl.DeleteFamilie_Bp)
	r.Get("/all", familieCtrl.GetAllFamilie_Bp)

	return r
}
