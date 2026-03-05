package models

import "time"

type Group struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	MemberCount     int       `json:"memberCount"`
	ExternalID      string    `json:"externalId,omitempty"`
	DisableMentions bool      `json:"disableMentions"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type GroupCreateParams struct {
	Name            string `json:"name"`
	ExternalID      string `json:"externalId,omitempty"`
	DisableMentions bool   `json:"disableMentions,omitempty"`
}

type GroupUpdateParams struct {
	ID              string `json:"id"`
	Name            string `json:"name,omitempty"`
	ExternalID      string `json:"externalId,omitempty"`
	DisableMentions *bool  `json:"disableMentions,omitempty"`
}

type GroupListParams struct {
	PaginationParams
	Sort       string `json:"sort,omitempty"`
	Direction  string `json:"direction,omitempty"`
	Query      string `json:"query,omitempty"`
	UserID     string `json:"userId,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
	Name       string `json:"name,omitempty"`
}

type GroupMembership struct {
	ID   string `json:"id"`
	User User   `json:"user"`
}
