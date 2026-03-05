package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type AttachmentsService struct {
	client *Client
}

func (s *AttachmentsService) List(ctx context.Context, params models.AttachmentListParams) ([]models.Attachment, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "attachments.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var attachments []models.Attachment
	if err := json.Unmarshal(resp.Data, &attachments); err != nil {
		return nil, nil, err
	}
	return attachments, resp.Pagination, nil
}

func (s *AttachmentsService) Create(ctx context.Context, params models.AttachmentCreateParams) (*models.Attachment, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "attachments.create", params, &resp); err != nil {
		return nil, err
	}
	var att models.Attachment
	if err := json.Unmarshal(resp.Data, &att); err != nil {
		return nil, err
	}
	return &att, nil
}

func (s *AttachmentsService) Delete(ctx context.Context, id string) error {
	return s.client.Post(ctx, "attachments.delete", map[string]string{"id": id}, nil)
}
