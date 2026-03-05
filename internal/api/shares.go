package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type SharesService struct {
	client *Client
}

func (s *SharesService) List(ctx context.Context, params models.ShareListParams) ([]models.Share, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "shares.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var shares []models.Share
	if err := json.Unmarshal(resp.Data, &shares); err != nil {
		return nil, nil, err
	}
	return shares, resp.Pagination, nil
}

func (s *SharesService) Create(ctx context.Context, params models.ShareCreateParams) (*models.Share, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "shares.create", params, &resp); err != nil {
		return nil, err
	}
	var share models.Share
	if err := json.Unmarshal(resp.Data, &share); err != nil {
		return nil, err
	}
	return &share, nil
}

func (s *SharesService) Update(ctx context.Context, params models.ShareUpdateParams) (*models.Share, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "shares.update", params, &resp); err != nil {
		return nil, err
	}
	var share models.Share
	if err := json.Unmarshal(resp.Data, &share); err != nil {
		return nil, err
	}
	return &share, nil
}

func (s *SharesService) Revoke(ctx context.Context, id string) error {
	return s.client.Post(ctx, "shares.revoke", map[string]string{"id": id}, nil)
}

func (s *SharesService) Info(ctx context.Context, id string) (*models.Share, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "shares.info", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var share models.Share
	if err := json.Unmarshal(resp.Data, &share); err != nil {
		return nil, err
	}
	return &share, nil
}
