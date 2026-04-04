package utils

import (
	"math"
	"net/http"
	"project-keuangan-keluarga/model"
	"slices"
	"strconv"
	"strings"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// ParsePaginationParams extracts and validates pagination query parameters
// from the incoming HTTP request. allowedSorts is the whitelist of column
// names the caller permits; defaultSort is used when the requested sort
// field is absent or not in the whitelist.
func ParsePaginationParams(r *http.Request, allowedSorts []string, defaultSort string) model.PaginationParams {
	q := r.URL.Query()

	// ── page ────────────────────────────────────────────────────
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	// ── limit ───────────────────────────────────────────────────
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	// ── sort ────────────────────────────────────────────────────
	sort := strings.TrimSpace(q.Get("sort"))
	if sort == "" || !slices.Contains(allowedSorts, sort) {
		sort = defaultSort
	}

	// ── order ───────────────────────────────────────────────────
	order := strings.ToLower(strings.TrimSpace(q.Get("order")))
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// ── search ──────────────────────────────────────────────────
	search := strings.TrimSpace(q.Get("search"))

	return model.PaginationParams{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
		Search: search,
	}
}

// CalculateOffset returns the SQL OFFSET value for a given page and limit.
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// BuildPaginationMeta constructs the PaginationMeta from raw totals.
func BuildPaginationMeta(totalItems, page, limit int) model.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return model.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}
