package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type RevisionsService struct {
	client *Client
}

func (s *RevisionsService) List(ctx context.Context, params models.RevisionListParams) ([]models.Revision, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "revisions.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var revisions []models.Revision
	if err := json.Unmarshal(resp.Data, &revisions); err != nil {
		return nil, nil, err
	}
	return revisions, resp.Pagination, nil
}

func (s *RevisionsService) Info(ctx context.Context, id string) (*models.Revision, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "revisions.info", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var rev models.Revision
	if err := json.Unmarshal(resp.Data, &rev); err != nil {
		return nil, err
	}
	return &rev, nil
}

func (s *RevisionsService) Delete(ctx context.Context, id string) error {
	return s.client.Post(ctx, "revisions.delete", map[string]string{"id": id}, nil)
}
