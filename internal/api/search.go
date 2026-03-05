package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type SearchService struct {
	client *Client
}

func (s *SearchService) Documents(ctx context.Context, params models.SearchParams) ([]models.SearchResult, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.search", params, &resp); err != nil {
		return nil, nil, err
	}
	var results []models.SearchResult
	if err := json.Unmarshal(resp.Data, &results); err != nil {
		return nil, nil, err
	}
	return results, resp.Pagination, nil
}

func (s *SearchService) Titles(ctx context.Context, params models.SearchParams) ([]models.Document, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.search_titles", params, &resp); err != nil {
		return nil, nil, err
	}
	var docs []models.Document
	if err := json.Unmarshal(resp.Data, &docs); err != nil {
		return nil, nil, err
	}
	return docs, resp.Pagination, nil
}
