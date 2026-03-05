package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type UsersService struct {
	client *Client
}

func (s *UsersService) List(ctx context.Context, params models.UserListParams) ([]models.User, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "users.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var users []models.User
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		return nil, nil, err
	}
	return users, resp.Pagination, nil
}

func (s *UsersService) Info(ctx context.Context, id string) (*models.User, error) {
	var resp models.APIResponse
	body := map[string]string{}
	if id != "" {
		body["id"] = id
	}
	if err := s.client.Post(ctx, "users.info", body, &resp); err != nil {
		return nil, err
	}
	var user models.User
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
