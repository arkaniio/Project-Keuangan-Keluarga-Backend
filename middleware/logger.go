package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ── Color codes for terminal output ─────────────────────────────────────────

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// ── Response Writer wrapper ─────────────────────────────────────────────────

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// newResponseWriter creates a wrapped ResponseWriter with 200 as default status.
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the status code before delegating to the original writer.
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

// Write ensures WriteHeader is called before writing the body.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// ── Log entry ───────────────────────────────────────────────────────────────

// LogEntry represents a single structured log record for an HTTP request.
type LogEntry struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	IP        string `json:"ip"`
	Time      string `json:"time"`
}

// ── Middleware ───────────────────────────────────────────────────────────────

// Logger returns an HTTP middleware that logs every request in structured JSON
// format and prints a colorized summary line to the terminal.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the writer to capture the status code
		wrapped := newResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Calculate latency
		latency := time.Since(start)
		latencyMs := latency.Milliseconds()

		// Resolve client IP
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		// Build structured log entry
		entry := LogEntry{
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    wrapped.statusCode,
			LatencyMs: latencyMs,
			IP:        ip,
			Time:      start.UTC().Format(time.RFC3339),
		}

		// JSON output (structured log)
		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			log.Printf("[LOGGER] Failed to marshal log entry: %v", err)
		} else {
			log.Println(string(jsonBytes))
		}

		// Colorized terminal line
		statusColor := colorForStatus(wrapped.statusCode)
		methodColor := colorForMethod(r.Method)

		fmt.Printf("%s[%s]%s %s%-7s%s %s → %s%d%s in %dms from %s\n",
			colorWhite, start.Format("15:04:05"), colorReset,
			methodColor, r.Method, colorReset,
			r.URL.Path,
			statusColor, wrapped.statusCode, colorReset,
			latencyMs,
			ip,
		)
	})
}

// ── Helpers ─────────────────────────────────────────────────────────────────

// colorForStatus returns a terminal color based on the HTTP status code range.
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen
	case code >= 300 && code < 400:
		return colorCyan
	case code >= 400 && code < 500:
		return colorYellow
	default:
		return colorRed
	}
}

// colorForMethod returns a terminal color based on the HTTP method.
func colorForMethod(method string) string {
	switch method {
	case http.MethodGet:
		return colorGreen
	case http.MethodPost:
		return colorCyan
	case http.MethodPut, http.MethodPatch:
		return colorYellow
	case http.MethodDelete:
		return colorRed
	default:
		return colorWhite
	}
}
