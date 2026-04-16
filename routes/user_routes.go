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

// UserRoutes creates the chi router for user-related endpoints.
// It accepts both a general limiter (for all routes) and a strict limiter
// (for auth-sensitive endpoints like login and register).
func UserRoutes(userCtrl *controller.ControllerHandler, generalLimiter *ratelimiter.Limiter, strictLimiter *ratelimiter.Limiter) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)                              // custom structured JSON logger
	r.Use(chimw.Recoverer)                                // recover from panics
	r.Use(chimw.RequestID)                                // inject X-Request-Id header
	r.Use(middleware.RateLimitMiddleware(generalLimiter)) // general rate limit: 60 req/min
	r.Use(middleware.CorsMiddleware())

	// Health-check endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// Auth routes with stricter rate limiting (10 req/min)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimitMiddleware(strictLimiter))
		r.Post("/register", userCtrl.Register)
		r.Post("/login", userCtrl.Login)
	})

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(middleware.MiddlewareAuth)
		r.Get("/profile", userCtrl.GetProfile)
		r.Put("/update", userCtrl.UpdateUser)
		r.Get("/all", userCtrl.GetAllUser)
	})

	return r
}
