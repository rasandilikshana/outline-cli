package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	LastActiveAt *time.Time `json:"lastActiveAt"`
	IsSuspended  bool      `json:"isSuspended"`
}

type UserListParams struct {
	PaginationParams
	Sort      string `json:"sort,omitempty"`
	Direction string `json:"direction,omitempty"`
	Query     string `json:"query,omitempty"`
	Role      string `json:"role,omitempty"`
	Filter    string `json:"filter,omitempty"`
}
