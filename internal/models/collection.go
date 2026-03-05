package models

import "time"

type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	Permission  string    `json:"permission"`
	Sharing     bool      `json:"sharing"`
	Commenting  bool      `json:"commenting"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CollectionCreateParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Permission  string `json:"permission,omitempty"`
	Sharing     *bool  `json:"sharing,omitempty"`
	Commenting  *bool  `json:"commenting,omitempty"`
}

type CollectionUpdateParams struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Permission  string `json:"permission,omitempty"`
	Sharing     *bool  `json:"sharing,omitempty"`
	Commenting  *bool  `json:"commenting,omitempty"`
}

type CollectionListParams struct {
	PaginationParams
	Query        string `json:"query,omitempty"`
	StatusFilter string `json:"statusFilter,omitempty"`
}
