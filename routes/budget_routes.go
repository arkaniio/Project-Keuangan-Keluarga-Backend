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
func BudgetRoutes(bgt_controller *controller.ControllerBudget, generalLimiter *ratelimiter.Limiter) *chi.Mux {
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
	r.Post("/", bgt_controller.CreateNewBudget_Bp)
	r.Put("/update", bgt_controller.UpdateBudget_Bp)
	r.Delete("/:id", bgt_controller.DeleteBudget_Bp)
	r.Get("/", bgt_controller.GetAllBudget_Bp)

	return r
}
