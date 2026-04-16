package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

// ─── CORS Configuration ────────────────────────────────────────────────────
//
// Industry best-practices applied:
//
//  1. Explicit origin whitelist — never use wildcard "*" in production.
//     Configurable via CORS_ALLOWED_ORIGINS env var (comma-separated).
//
//  2. Restricted HTTP methods — only methods the API actually uses.
//
//  3. Minimal allowed headers — only headers the frontend needs to send.
//     Authorization (JWT), Content-Type (JSON), and X-Request-Id.
//
//  4. Exposed headers — lets the frontend read rate-limit and request-id
//     headers from responses.
//
//  5. Credentials support — enabled so the browser can send cookies and
//     Authorization headers cross-origin.
//
//  6. Max age (preflight cache) — browsers cache preflight (OPTIONS) for
//     5 minutes, reducing unnecessary round-trips.
//
//  7. OPTIONS passthrough disabled — chi/cors handles OPTIONS internally
//     to prevent unprotected preflight leaks to handlers.
//
// ────────────────────────────────────────────────────────────────────────────

// CorsMiddleware returns a production-ready CORS handler.
//
// Environment variables:
//
//	CORS_ALLOWED_ORIGINS  — comma-separated list of allowed origins.
//	                        Default: "http://localhost:3000,http://localhost:5173"
//	                        (common dev ports for Next.js & Vite)
func CorsMiddleware() func(http.Handler) http.Handler {

	// ── 1. Parse allowed origins from env ────────────────────────
	allowedOrigins := []string{
		"http://localhost:3000", // Next.js default
		"http://localhost:5173", // Vite default
	}

	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		origins := strings.Split(envOrigins, ",")
		allowedOrigins = make([]string, 0, len(origins))
		for _, o := range origins {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	log.Printf("[CORS] Allowed origins: %v", allowedOrigins)

	// ── 2. Build CORS options ───────────────────────────────────
	return cors.Handler(cors.Options{

		// Only allow requests from explicitly whitelisted origins.
		// Never use "*" — it disables credential support and opens
		// the API to any website.
		AllowedOrigins: allowedOrigins,

		// Only the HTTP methods this API actually serves.
		// PATCH is included for future partial-update endpoints.
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},

		// Headers the client is allowed to send.
		//  - Authorization  : JWT Bearer tokens
		//  - Content-Type   : application/json, multipart/form-data
		//  - Accept         : standard content negotiation
		//  - X-Request-Id   : request tracing (chi middleware)
		//  - X-CSRF-Token   : CSRF protection (if added later)
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
			"X-Request-Id",
			"X-CSRF-Token",
		},

		// Headers the browser is allowed to read from the response.
		// By default browsers block all non-simple response headers
		// in cross-origin requests.
		ExposedHeaders: []string{
			"X-Request-Id",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},

		// Allow credentials (cookies, Authorization header).
		// Required for JWT-based auth from a cross-origin frontend.
		AllowCredentials: true,

		// How long (seconds) the browser should cache the preflight
		// response. 300s = 5 minutes is a good balance between
		// reducing OPTIONS requests and respecting config changes.
		MaxAge: 300,
	})
}
