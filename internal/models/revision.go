package models

import "time"

type Revision struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"documentId"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  *User     `json:"createdBy,omitempty"`
}

type RevisionListParams struct {
	PaginationParams
	DocumentID string `json:"documentId"`
	Sort       string `json:"sort,omitempty"`
	Direction  string `json:"direction,omitempty"`
}
