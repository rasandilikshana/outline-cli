package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type CommentsService struct {
	client *Client
}

func (s *CommentsService) List(ctx context.Context, params models.CommentListParams) ([]models.Comment, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "comments.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var comments []models.Comment
	if err := json.Unmarshal(resp.Data, &comments); err != nil {
		return nil, nil, err
	}
	return comments, resp.Pagination, nil
}

func (s *CommentsService) Create(ctx context.Context, params models.CommentCreateParams) (*models.Comment, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "comments.create", params, &resp); err != nil {
		return nil, err
	}
	var comment models.Comment
	if err := json.Unmarshal(resp.Data, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentsService) Update(ctx context.Context, params models.CommentUpdateParams) (*models.Comment, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "comments.update", params, &resp); err != nil {
		return nil, err
	}
	var comment models.Comment
	if err := json.Unmarshal(resp.Data, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentsService) Delete(ctx context.Context, id string) error {
	return s.client.Post(ctx, "comments.delete", map[string]string{"id": id}, nil)
}

func (s *CommentsService) Resolve(ctx context.Context, id string) (*models.Comment, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "comments.resolve", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var comment models.Comment
	if err := json.Unmarshal(resp.Data, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentsService) Unresolve(ctx context.Context, id string) (*models.Comment, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "comments.unresolve", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var comment models.Comment
	if err := json.Unmarshal(resp.Data, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}
