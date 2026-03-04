package common

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
)

// PaginationRequest represents cursor-based pagination input.
type PaginationRequest struct {
	Cursor string `form:"cursor"`
	Limit  int    `form:"limit"`
}

// PaginationResponse represents cursor-based pagination output,
// matching the API contract format.
type PaginationResponse struct {
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
	TotalCount int64  `json:"total_count,omitempty"`
}

// DefaultLimit is the default page size when none is specified.
const DefaultLimit = 20

// MaxLimit is the maximum allowed page size.
const MaxLimit = 100

// EffectiveLimit returns the limit, clamped to the allowed range.
func (p *PaginationRequest) EffectiveLimit() int {
	if p.Limit <= 0 {
		return DefaultLimit
	}
	if p.Limit > MaxLimit {
		return MaxLimit
	}
	return p.Limit
}

// CursorData holds the decoded cursor state.
type CursorData struct {
	ID        string `json:"id"`
	Timestamp string `json:"ts,omitempty"`
}

// DecodeCursor decodes a base64-encoded cursor string.
func DecodeCursor(cursor string) (*CursorData, error) {
	if cursor == "" {
		return nil, nil
	}
	decoded, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	var data CursorData
	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// EncodeCursor encodes cursor data to a base64 string.
func EncodeCursor(id string, timestamp string) string {
	data := CursorData{ID: id, Timestamp: timestamp}
	b, _ := json.Marshal(data)
	return base64.URLEncoding.EncodeToString(b)
}

// ParseIntParam safely parses a string to int with a default value.
func ParseIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
