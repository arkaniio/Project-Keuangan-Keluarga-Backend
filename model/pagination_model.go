package model

// PaginationParams holds the parsed pagination parameters from query string.
type PaginationParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// PaginationMeta holds the pagination metadata included in API responses.
type PaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

// PaginatedResponse wraps a list of items together with pagination metadata.
type PaginatedResponse struct {
	Items      interface{}    `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}
