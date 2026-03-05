package models

import "time"

type Event struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	ModelID      string    `json:"modelId"`
	ActorID      string    `json:"actorId"`
	DocumentID   string    `json:"documentId,omitempty"`
	CollectionID string    `json:"collectionId,omitempty"`
	IP           string    `json:"ip,omitempty"`
	Data         any       `json:"data,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	Actor        *User     `json:"actor,omitempty"`
}

type EventListParams struct {
	PaginationParams
	Sort         string   `json:"sort,omitempty"`
	Direction    string   `json:"direction,omitempty"`
	Name         string   `json:"name,omitempty"`
	Events       []string `json:"events,omitempty"`
	AuditLog     bool     `json:"auditLog,omitempty"`
	ActorID      string   `json:"actorId,omitempty"`
	DocumentID   string   `json:"documentId,omitempty"`
	CollectionID string   `json:"collectionId,omitempty"`
}
