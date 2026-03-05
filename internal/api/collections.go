package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type CollectionsService struct {
	client *Client
}

func (s *CollectionsService) List(ctx context.Context, params models.CollectionListParams) ([]models.Collection, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var collections []models.Collection
	if err := json.Unmarshal(resp.Data, &collections); err != nil {
		return nil, nil, err
	}
	return collections, resp.Pagination, nil
}

func (s *CollectionsService) Info(ctx context.Context, id string) (*models.Collection, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.info", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var coll models.Collection
	if err := json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}
	return &coll, nil
}

func (s *CollectionsService) Create(ctx context.Context, params models.CollectionCreateParams) (*models.Collection, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.create", params, &resp); err != nil {
		return nil, err
	}
	var coll models.Collection
	if err := json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}
	return &coll, nil
}

func (s *CollectionsService) Update(ctx context.Context, params models.CollectionUpdateParams) (*models.Collection, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.update", params, &resp); err != nil {
		return nil, err
	}
	var coll models.Collection
	if err := json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}
	return &coll, nil
}

func (s *CollectionsService) Delete(ctx context.Context, id string) error {
	return s.client.Post(ctx, "collections.delete", map[string]string{"id": id}, nil)
}

func (s *CollectionsService) Archive(ctx context.Context, id string) (*models.Collection, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.archive", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var coll models.Collection
	if err := json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}
	return &coll, nil
}

func (s *CollectionsService) Restore(ctx context.Context, id string) (*models.Collection, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.restore", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var coll models.Collection
	if err := json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}
	return &coll, nil
}

func (s *CollectionsService) Documents(ctx context.Context, id string) ([]models.NavigationNode, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "collections.documents", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var nodes []models.NavigationNode
	if err := json.Unmarshal(resp.Data, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (s *CollectionsService) FindByName(ctx context.Context, name string) (*models.Collection, error) {
	collections, _, err := s.List(ctx, models.CollectionListParams{
		PaginationParams: models.PaginationParams{Limit: 100},
	})
	if err != nil {
		return nil, err
	}
	for _, c := range collections {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, nil
}
