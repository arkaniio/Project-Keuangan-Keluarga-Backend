package ratelimiter

import (
	"math"
	"time"
)

// ── Rate Limit Info ─────────────────────────────────────────────────────────

// RateLimitInfo contains rate limit state to be communicated back to clients
// via HTTP response headers.
type RateLimitInfo struct {
	// Limit is the maximum number of requests allowed (Rate + BurstCapacity).
	Limit int

	// Remaining is the number of requests the client can still make
	// within the current window.
	Remaining int

	// RetryAfterSec is the number of seconds the client should wait
	// before retrying. Only meaningful when the request is denied.
	RetryAfterSec int
}

// ── Limiter ─────────────────────────────────────────────────────────────────

// Limiter implements the Sliding Window Counter rate limiting algorithm.
//
// How it works:
//  1. Time is divided into fixed windows (e.g., 1 minute each).
//  2. For each key (IP or user), we track request counts in the current and
//     previous windows.
//  3. The effective count is calculated as:
//     weighted = prevCount × (1 - elapsedFraction) + currCount
//  4. If weighted >= limit (Rate + BurstCapacity), the request is denied.
//
// This produces a smooth, accurate rate limit without the boundary-spike
// problem of fixed windows, while using far less memory than a full
// sliding window log.
type Limiter struct {
	config Config
	store  Store
}

// NewLimiter creates a new Limiter with the given config and store.
func NewLimiter(cfg Config, store Store) *Limiter {
	return &Limiter{
		config: cfg,
		store:  store,
	}
}

// Allow checks whether a request from the given key should be allowed.
// Returns whether the request is allowed and rate limit info for headers.
func (l *Limiter) Allow(key string) (bool, RateLimitInfo) {
	effectiveLimit := l.config.Rate + l.config.BurstCapacity

	prevCount, currCount, windowStart, err := l.store.Increment(key, l.config.Window)
	if err != nil {
		// On store error, fail open (allow the request) to avoid
		// blocking legitimate traffic due to infrastructure issues.
		return true, RateLimitInfo{
			Limit:     effectiveLimit,
			Remaining: effectiveLimit,
		}
	}

	// ── Sliding Window Counter calculation ───────────────────────────────
	//
	// elapsed    = time since current window started
	// fraction   = elapsed / windowSize  (0.0 → 1.0)
	// prevWeight = 1.0 - fraction
	// weighted   = prevCount × prevWeight + currCount
	//
	// Example: window=60s, elapsed=20s, prevCount=50, currCount=10
	//   fraction   = 20/60 = 0.333
	//   prevWeight = 0.667
	//   weighted   = 50×0.667 + 10 = 43.33
	//
	now := time.Now()
	elapsed := now.Sub(windowStart)
	windowDuration := l.config.Window

	var fraction float64
	if windowDuration > 0 {
		fraction = float64(elapsed) / float64(windowDuration)
		if fraction > 1.0 {
			fraction = 1.0
		}
	}

	prevWeight := 1.0 - fraction
	weightedCount := float64(prevCount)*prevWeight + float64(currCount)

	remaining := effectiveLimit - int(math.Ceil(weightedCount))
	if remaining < 0 {
		remaining = 0
	}

	if weightedCount > float64(effectiveLimit) {
		// Request denied — calculate retry-after
		// Estimate how long until enough of the previous window's weight
		// has decayed to bring the count below the limit.
		retryAfter := int(math.Ceil(windowDuration.Seconds() - elapsed.Seconds()))
		if retryAfter < 1 {
			retryAfter = 1
		}

		return false, RateLimitInfo{
			Limit:         effectiveLimit,
			Remaining:     0,
			RetryAfterSec: retryAfter,
		}
	}

	return true, RateLimitInfo{
		Limit:     effectiveLimit,
		Remaining: remaining,
	}
}

// IsWhitelisted checks if the given IP is in the whitelist.
func (l *Limiter) IsWhitelisted(ip string) bool {
	for _, whitelistedIP := range l.config.IPWhitelist {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}

// GetConfig returns a copy of the limiter's configuration.
func (l *Limiter) GetConfig() Config {
	return l.config
}
