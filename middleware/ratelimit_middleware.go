package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"project-keuangan-keluarga/middleware/ratelimiter"
	"project-keuangan-keluarga/utils"

	"github.com/google/uuid"
)

// ── Rate Limit Middleware ───────────────────────────────────────────────────

// RateLimitMiddleware returns a Chi-compatible middleware that enforces
// rate limiting using the provided Limiter instance.
//
// Key selection strategy:
//   - If the request context contains an authenticated user ID (set by
//     MiddlewareAuth), the key is "user:<uuid>" for per-user limiting.
//   - Otherwise, the key is "ip:<client-ip>" for per-IP limiting.
//
// Whitelisted IPs bypass rate limiting entirely.
//
// On rate limit exceeded:
//   - Responds with 429 Too Many Requests
//   - Sets X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After headers
//   - Logs the blocked request for monitoring
func RateLimitMiddleware(limiter *ratelimiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ── 1. Extract client IP ────────────────────────────────
			clientIP := extractClientIP(r)

			// ── 2. Check whitelist ──────────────────────────────────
			if limiter.IsWhitelisted(clientIP) {
				next.ServeHTTP(w, r)
				return
			}

			// ── 3. Build rate limit key ─────────────────────────────
			key := buildRateLimitKey(r, clientIP)

			// ── 4. Check rate limit ─────────────────────────────────
			allowed, info := limiter.Allow(key)

			// ── 5. Set rate limit headers (always) ──────────────────
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(info.Limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(info.Remaining))

			if !allowed {
				// ── 6a. Rate limit exceeded ─────────────────────────
				retryAfter := strconv.Itoa(info.RetryAfterSec)
				w.Header().Set("Retry-After", retryAfter)

				// Log the blocked request
				log.Printf(
					"[RATE-LIMIT] Blocked request: ip=%s method=%s path=%s key=%s retry_after=%ss",
					clientIP, r.Method, r.URL.Path, key, retryAfter,
				)

				// Colorized terminal warning
				fmt.Printf(
					"%s[RATE-LIMIT]%s %s%s%s %s → %s429 Too Many Requests%s (key=%s, retry=%ss)\n",
					colorYellow, colorReset,
					colorForMethod(r.Method), r.Method, colorReset,
					r.URL.Path,
					colorRed, colorReset,
					key, retryAfter,
				)

				utils.ResponseError(
					w,
					http.StatusTooManyRequests,
					"Too many requests. Please try again later.",
					map[string]interface{}{
						"retry_after_seconds": info.RetryAfterSec,
						"limit":              info.Limit,
					},
				)
				return
			}

			// ── 6b. Request allowed ─────────────────────────────────
			next.ServeHTTP(w, r)
		})
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────────

// extractClientIP resolves the real client IP from the request,
// checking proxy headers in order of trust.
func extractClientIP(r *http.Request) string {
	// X-Forwarded-For may contain a chain: "client, proxy1, proxy2"
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}

	// X-Real-IP is set by some reverse proxies (e.g., Nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr (strip port)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// buildRateLimitKey creates the appropriate rate limit key.
// Uses authenticated user ID if available, otherwise falls back to IP.
func buildRateLimitKey(r *http.Request, clientIP string) string {
	// Check if user is authenticated (set by MiddlewareAuth)
	if userID, ok := r.Context().Value("id").(uuid.UUID); ok && userID != uuid.Nil {
		return "user:" + userID.String()
	}
	return "ip:" + clientIP
}
