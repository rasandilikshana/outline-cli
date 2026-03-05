package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type AuthService struct {
	client *Client
}

type AuthInfo struct {
	User models.User `json:"user"`
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
}

func (s *AuthService) Info(ctx context.Context) (*AuthInfo, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "auth.info", nil, &resp); err != nil {
		return nil, err
	}
	var info AuthInfo
	if err := json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
