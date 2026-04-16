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

func FamilyMemberRoutes(familyMemberCtrl controller.ControllerHandlerFamilyMember, generalLimiter *ratelimiter.Limiter) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)                              // custom structured JSON logger
	r.Use(middleware.MiddlewareAuth)                      // auth required
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

	// API v1 routes
	r.Post("/", familyMemberCtrl.CreateFamilyMember_Bp)
	r.Put("/update", familyMemberCtrl.UpdateFamilyMember_Bp)
	r.Delete("/delete", familyMemberCtrl.DeleteFamilyMember_Bp)
	r.Get("/all", familyMemberCtrl.GetAllFamilyMember_Bp)
	r.Get("/me", familyMemberCtrl.GetMyMembership_Bp)

	return r
}
