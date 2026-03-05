package models

import "time"

type Comment struct {
	ID              string    `json:"id"`
	DocumentID      string    `json:"documentId"`
	ParentCommentID *string   `json:"parentCommentId"`
	Data            any       `json:"data"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedBy       *User     `json:"createdBy,omitempty"`
}

type CommentCreateParams struct {
	DocumentID      string `json:"documentId"`
	ParentCommentID string `json:"parentCommentId,omitempty"`
	Text            string `json:"text,omitempty"`
	Data            any    `json:"data,omitempty"`
}

type CommentUpdateParams struct {
	ID   string `json:"id"`
	Data any    `json:"data,omitempty"`
}

type CommentListParams struct {
	PaginationParams
	Sort            string `json:"sort,omitempty"`
	Direction       string `json:"direction,omitempty"`
	DocumentID      string `json:"documentId,omitempty"`
	ParentCommentID string `json:"parentCommentId,omitempty"`
	CollectionID    string `json:"collectionId,omitempty"`
	StatusFilter    string `json:"statusFilter,omitempty"`
}
