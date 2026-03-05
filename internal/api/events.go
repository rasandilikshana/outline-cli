package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type EventsService struct {
	client *Client
}

func (s *EventsService) List(ctx context.Context, params models.EventListParams) ([]models.Event, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "events.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var events []models.Event
	if err := json.Unmarshal(resp.Data, &events); err != nil {
		return nil, nil, err
	}
	return events, resp.Pagination, nil
}
