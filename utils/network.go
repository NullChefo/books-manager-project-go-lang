package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// PaginationParams holds the pagination parameters
type PaginationParams struct {
	Page int
	Size int
}

// PaginatedResponse holds the paginated response structure
type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// ParsePaginationParams parses the pagination parameters from the request
func ParsePaginationParams(r *http.Request) PaginationParams {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	page := 1
	size := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = s
		}
	}

	return PaginationParams{
		Page: page,
		Size: size,
	}
}

// GetPaginatedData fetches paginated data from the database
func GetPaginatedData[T any](offset, size int, fetchFunc func(int, int) ([]T, error)) ([]T, error) {
	return fetchFunc(offset, size)
}

// CreatePaginatedResponse creates a paginated response structure
func CreatePaginatedResponse[T any](items []T, page, size, totalItems int) PaginatedResponse[T] {
	totalPages := (totalItems + size - 1) / size
	return PaginatedResponse[T]{
		Items:      items,
		Page:       page,
		Size:       size,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// ParseID extracts and parses the ID from the URL path.
func ParseID(r *http.Request) (int64, error) {
	// Extract the book ID from the URL path
	idStr := r.PathValue("id")
	return strconv.ParseInt(idStr, 10, 64)
}

// SendJSONResponse sends a JSON response with the given data and status code.
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
