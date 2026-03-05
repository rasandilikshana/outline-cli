package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type StarsService struct {
	client *Client
}

func (s *StarsService) List(ctx context.Context, params models.StarListParams) ([]models.Star, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "stars.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var stars []models.Star
	if err := json.Unmarshal(resp.Data, &stars); err != nil {
		return nil, nil, err
	}
	return stars, resp.Pagination, nil
}

func (s *StarsService) Create(ctx context.Context, params models.StarCreateParams) (*models.Star, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "stars.create", params, &resp); err != nil {
		return nil, err
	}
	var star models.Star
	if err := json.Unmarshal(resp.Data, &star); err != nil {
		return nil, err
	}
	return &star, nil
}

func (s *StarsService) Delete(ctx context.Context, id string) error {
	return s.client.Post(ctx, "stars.delete", map[string]string{"id": id}, nil)
}
