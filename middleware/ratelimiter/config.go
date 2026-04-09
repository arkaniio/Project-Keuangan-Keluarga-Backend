package ratelimiter

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// ── Config ──────────────────────────────────────────────────────────────────

// Config holds all rate limiter settings.
type Config struct {
	// Rate is the maximum number of requests allowed per window.
	Rate int

	// Window is the time window for counting requests.
	Window time.Duration

	// BurstCapacity is the number of extra requests allowed above the rate
	// to accommodate short bursts. The effective limit = Rate + BurstCapacity.
	BurstCapacity int

	// IPWhitelist is a list of IPs that bypass rate limiting entirely.
	IPWhitelist []string
}

// ── Defaults ────────────────────────────────────────────────────────────────

// DefaultConfig returns a general-purpose rate limit configuration.
// 60 requests per minute with burst capacity of 10.
func DefaultConfig() Config {
	return Config{
		Rate:          60,
		Window:        1 * time.Minute,
		BurstCapacity: 10,
		IPWhitelist:   []string{"127.0.0.1", "::1"},
	}
}

// StrictConfig returns a stricter configuration for sensitive endpoints
// such as login and registration. 10 requests per minute with burst of 3.
func StrictConfig() Config {
	return Config{
		Rate:          10,
		Window:        1 * time.Minute,
		BurstCapacity: 3,
		IPWhitelist:   []string{"127.0.0.1", "::1"},
	}
}

// ── Environment Loader ──────────────────────────────────────────────────────

// LoadFromEnv reads rate limit configuration from environment variables.
// Falls back to DefaultConfig values if variables are not set.
//
// Environment variables:
//   - RATE_LIMIT_RATE           → requests per window (default: 60)
//   - RATE_LIMIT_WINDOW_SECONDS → window duration in seconds (default: 60)
//   - RATE_LIMIT_BURST          → burst capacity (default: 10)
//   - RATE_LIMIT_WHITELIST      → comma-separated IPs (default: "127.0.0.1,::1")
func LoadFromEnv() Config {
	cfg := DefaultConfig()

	if v := os.Getenv("RATE_LIMIT_RATE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Rate = n
		}
	}

	if v := os.Getenv("RATE_LIMIT_WINDOW_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Window = time.Duration(n) * time.Second
		}
	}

	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.BurstCapacity = n
		}
	}

	if v := os.Getenv("RATE_LIMIT_WHITELIST"); v != "" {
		ips := strings.Split(v, ",")
		cleaned := make([]string, 0, len(ips))
		for _, ip := range ips {
			if trimmed := strings.TrimSpace(ip); trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		cfg.IPWhitelist = cleaned
	}

	return cfg
}

// LoadStrictFromEnv reads strict rate limit configuration from environment
// variables. Falls back to StrictConfig values if variables are not set.
//
// Environment variables:
//   - RATE_LIMIT_STRICT_RATE           → requests per window (default: 10)
//   - RATE_LIMIT_STRICT_WINDOW_SECONDS → window duration in seconds (default: 60)
//   - RATE_LIMIT_STRICT_BURST          → burst capacity (default: 3)
func LoadStrictFromEnv() Config {
	cfg := StrictConfig()

	if v := os.Getenv("RATE_LIMIT_STRICT_RATE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Rate = n
		}
	}

	if v := os.Getenv("RATE_LIMIT_STRICT_WINDOW_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Window = time.Duration(n) * time.Second
		}
	}

	if v := os.Getenv("RATE_LIMIT_STRICT_BURST"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.BurstCapacity = n
		}
	}

	// Inherit whitelist from general config
	if v := os.Getenv("RATE_LIMIT_WHITELIST"); v != "" {
		ips := strings.Split(v, ",")
		cleaned := make([]string, 0, len(ips))
		for _, ip := range ips {
			if trimmed := strings.TrimSpace(ip); trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		cfg.IPWhitelist = cleaned
	}

	return cfg
}
