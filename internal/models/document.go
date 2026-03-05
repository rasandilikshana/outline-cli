package models

import "time"

type Document struct {
	ID               string     `json:"id"`
	CollectionID     string     `json:"collectionId"`
	ParentDocumentID *string    `json:"parentDocumentId"`
	Title            string     `json:"title"`
	Text             string     `json:"text"`
	Icon             string     `json:"icon,omitempty"`
	Color            string     `json:"color,omitempty"`
	Template         bool       `json:"template"`
	FullWidth        bool       `json:"fullWidth"`
	Revision         int        `json:"revision"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	PublishedAt      *time.Time `json:"publishedAt"`
	ArchivedAt       *time.Time `json:"archivedAt"`
	DeletedAt        *time.Time `json:"deletedAt"`
	CreatedBy        *User      `json:"createdBy,omitempty"`
	UpdatedBy        *User      `json:"updatedBy,omitempty"`
}

type NavigationNode struct {
	ID       string           `json:"id"`
	Title    string           `json:"title"`
	URL      string           `json:"url"`
	Children []NavigationNode `json:"children"`
}

type DocumentCreateParams struct {
	Title            string  `json:"title"`
	Text             string  `json:"text,omitempty"`
	CollectionID     string  `json:"collectionId,omitempty"`
	ParentDocumentID string  `json:"parentDocumentId,omitempty"`
	Publish          bool    `json:"publish"`
	Icon             string  `json:"icon,omitempty"`
	Color            string  `json:"color,omitempty"`
	FullWidth        bool    `json:"fullWidth,omitempty"`
	TemplateID       string  `json:"templateId,omitempty"`
}

type DocumentUpdateParams struct {
	ID       string `json:"id"`
	Title    string `json:"title,omitempty"`
	Text     string `json:"text,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Color    string `json:"color,omitempty"`
	Publish  *bool  `json:"publish,omitempty"`
	Append   bool   `json:"append,omitempty"`
}

type DocumentListParams struct {
	PaginationParams
	Sort             string `json:"sort,omitempty"`
	Direction        string `json:"direction,omitempty"`
	CollectionID     string `json:"collectionId,omitempty"`
	ParentDocumentID string `json:"parentDocumentId,omitempty"`
	UserID           string `json:"userId,omitempty"`
	StatusFilter     string `json:"statusFilter,omitempty"`
}

type DocumentMoveParams struct {
	ID               string `json:"id"`
	CollectionID     string `json:"collectionId,omitempty"`
	ParentDocumentID string `json:"parentDocumentId,omitempty"`
	Index            int    `json:"index,omitempty"`
}

type SearchParams struct {
	Query        string `json:"query"`
	CollectionID string `json:"collectionId,omitempty"`
	UserID       string `json:"userId,omitempty"`
	DateFilter   string `json:"dateFilter,omitempty"`
	PaginationParams
}

type SearchResult struct {
	Ranking  float64  `json:"ranking"`
	Context  string   `json:"context"`
	Document Document `json:"document"`
}
