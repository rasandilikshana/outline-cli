package models

import "encoding/json"

type APIResponse struct {
	OK         bool             `json:"ok"`
	Data       json.RawMessage  `json:"data,omitempty"`
	Error      string           `json:"error,omitempty"`
	Message    string           `json:"message,omitempty"`
	Status     int              `json:"status,omitempty"`
	Pagination *Pagination      `json:"pagination,omitempty"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type PaginationParams struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}
