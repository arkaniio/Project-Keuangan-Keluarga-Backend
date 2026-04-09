package ratelimiter

import (
	"sync"
	"time"
)

// ── Store Interface ─────────────────────────────────────────────────────────
// Store defines the contract for rate limit state persistence.
// Implement this interface to swap in Redis, Memcached, etc.

type Store interface {
	// Increment records a hit for the given key and returns the counts
	// for the previous and current windows, plus the current window start time.
	Increment(key string, window time.Duration) (prevCount int64, currCount int64, windowStart time.Time, err error)
}

// ── In-Memory Store ─────────────────────────────────────────────────────────

// windowEntry holds request counts for two adjacent fixed windows,
// used by the sliding window counter algorithm.
type windowEntry struct {
	prevCount   int64
	currCount   int64
	windowStart time.Time
}

// MemoryStore is a thread-safe, in-memory implementation of Store.
// Suitable for single-instance deployments.
type MemoryStore struct {
	mu      sync.RWMutex
	entries map[string]*windowEntry
	stopCh  chan struct{}
}

// NewMemoryStore creates a new MemoryStore and starts a background
// goroutine that periodically evicts expired entries.
func NewMemoryStore(cleanupInterval time.Duration) *MemoryStore {
	ms := &MemoryStore{
		entries: make(map[string]*windowEntry),
		stopCh:  make(chan struct{}),
	}
	go ms.cleanup(cleanupInterval)
	return ms
}

// Increment records a request hit and returns window counts.
// It automatically rotates windows when the current window expires.
func (ms *MemoryStore) Increment(key string, window time.Duration) (prevCount int64, currCount int64, windowStart time.Time, err error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()
	entry, exists := ms.entries[key]

	if !exists {
		// First request for this key
		ws := now.Truncate(window)
		entry = &windowEntry{
			prevCount:   0,
			currCount:   1,
			windowStart: ws,
		}
		ms.entries[key] = entry
		return 0, 1, ws, nil
	}

	// Calculate which window we're in
	currentWindowStart := now.Truncate(window)

	if currentWindowStart.After(entry.windowStart.Add(window)) {
		// We've skipped at least one full window — reset everything
		entry.prevCount = 0
		entry.currCount = 1
		entry.windowStart = currentWindowStart
	} else if currentWindowStart.After(entry.windowStart) {
		// We've moved to the next window — rotate
		entry.prevCount = entry.currCount
		entry.currCount = 1
		entry.windowStart = currentWindowStart
	} else {
		// Still in the same window
		entry.currCount++
	}

	return entry.prevCount, entry.currCount, entry.windowStart, nil
}

// cleanup periodically removes expired entries to prevent memory leaks.
func (ms *MemoryStore) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.evictExpired(interval)
		case <-ms.stopCh:
			return
		}
	}
}

// evictExpired removes entries whose window has not been updated
// within the given threshold.
func (ms *MemoryStore) evictExpired(maxAge time.Duration) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	cutoff := time.Now().Add(-2 * maxAge)
	for key, entry := range ms.entries {
		if entry.windowStart.Before(cutoff) {
			delete(ms.entries, key)
		}
	}
}

// Stop signals the cleanup goroutine to exit.
// Call this during graceful shutdown.
func (ms *MemoryStore) Stop() {
	close(ms.stopCh)
}
