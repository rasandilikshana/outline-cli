package models

import "time"

type Attachment struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ContentType string    `json:"contentType"`
	Size        int64     `json:"size"`
	URL         string    `json:"url"`
	DocumentID  string    `json:"documentId,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AttachmentCreateParams struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
	DocumentID  string `json:"documentId,omitempty"`
}

type AttachmentListParams struct {
	PaginationParams
	DocumentID string `json:"documentId,omitempty"`
	UserID     string `json:"userId,omitempty"`
}
