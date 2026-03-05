package models

import "time"

type Star struct {
	ID           string    `json:"id"`
	DocumentID   string    `json:"documentId,omitempty"`
	CollectionID string    `json:"collectionId,omitempty"`
	Index        string    `json:"index"`
	CreatedAt    time.Time `json:"createdAt"`
}

type StarCreateParams struct {
	DocumentID   string `json:"documentId,omitempty"`
	CollectionID string `json:"collectionId,omitempty"`
	Index        string `json:"index,omitempty"`
}

type StarListParams struct {
	PaginationParams
}
