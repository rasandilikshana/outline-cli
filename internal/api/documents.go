package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type DocumentsService struct {
	client *Client
}

func (s *DocumentsService) List(ctx context.Context, params models.DocumentListParams) ([]models.Document, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var docs []models.Document
	if err := json.Unmarshal(resp.Data, &docs); err != nil {
		return nil, nil, err
	}
	return docs, resp.Pagination, nil
}

func (s *DocumentsService) Info(ctx context.Context, id string) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.info", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Create(ctx context.Context, params models.DocumentCreateParams) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.create", params, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Update(ctx context.Context, params models.DocumentUpdateParams) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.update", params, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Delete(ctx context.Context, id string, permanent bool) error {
	return s.client.Post(ctx, "documents.delete", map[string]interface{}{"id": id, "permanent": permanent}, nil)
}

func (s *DocumentsService) Archive(ctx context.Context, id string) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.archive", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Restore(ctx context.Context, id string) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.restore", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Move(ctx context.Context, params models.DocumentMoveParams) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.move", params, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Export(ctx context.Context, id string) (string, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.export", map[string]string{"id": id}, &resp); err != nil {
		return "", err
	}
	var data string
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}
	return data, nil
}

func (s *DocumentsService) Duplicate(ctx context.Context, id string, recursive bool) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.duplicate", map[string]interface{}{"id": id, "recursive": recursive}, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentsService) Search(ctx context.Context, params models.SearchParams) ([]models.SearchResult, *models.Pagination, error) {
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

func (s *DocumentsService) Drafts(ctx context.Context, params models.DocumentListParams) ([]models.Document, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.drafts", params, &resp); err != nil {
		return nil, nil, err
	}
	var docs []models.Document
	if err := json.Unmarshal(resp.Data, &docs); err != nil {
		return nil, nil, err
	}
	return docs, resp.Pagination, nil
}

func (s *DocumentsService) Viewed(ctx context.Context, params models.PaginationParams) ([]models.Document, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.viewed", params, &resp); err != nil {
		return nil, nil, err
	}
	var docs []models.Document
	if err := json.Unmarshal(resp.Data, &docs); err != nil {
		return nil, nil, err
	}
	return docs, resp.Pagination, nil
}

func (s *DocumentsService) Unpublish(ctx context.Context, id string) (*models.Document, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "documents.unpublish", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var doc models.Document
	if err := json.Unmarshal(resp.Data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}
