package models

import "time"

type Share struct {
	ID                    string    `json:"id"`
	DocumentID            string    `json:"documentId,omitempty"`
	CollectionID          string    `json:"collectionId,omitempty"`
	Published             bool      `json:"published"`
	IncludeChildDocuments bool      `json:"includeChildDocuments"`
	URL                   string    `json:"url"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	CreatedBy             *User     `json:"createdBy,omitempty"`
}

type ShareCreateParams struct {
	DocumentID            string `json:"documentId,omitempty"`
	CollectionID          string `json:"collectionId,omitempty"`
	Published             *bool  `json:"published,omitempty"`
	IncludeChildDocuments *bool  `json:"includeChildDocuments,omitempty"`
}

type ShareUpdateParams struct {
	ID                    string `json:"id"`
	Published             *bool  `json:"published,omitempty"`
	IncludeChildDocuments *bool  `json:"includeChildDocuments,omitempty"`
}

type ShareListParams struct {
	PaginationParams
	Sort      string `json:"sort,omitempty"`
	Direction string `json:"direction,omitempty"`
	Query     string `json:"query,omitempty"`
}
