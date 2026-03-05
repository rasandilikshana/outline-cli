package api

import (
	"context"
	"encoding/json"
	"outline-cli/internal/models"
)

type GroupsService struct {
	client *Client
}

func (s *GroupsService) List(ctx context.Context, params models.GroupListParams) ([]models.Group, *models.Pagination, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "groups.list", params, &resp); err != nil {
		return nil, nil, err
	}
	var groups []models.Group
	if err := json.Unmarshal(resp.Data, &groups); err != nil {
		return nil, nil, err
	}
	return groups, resp.Pagination, nil
}

func (s *GroupsService) Info(ctx context.Context, id string) (*models.Group, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "groups.info", map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	var group models.Group
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupsService) Create(ctx context.Context, params models.GroupCreateParams) (*models.Group, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "groups.create", params, &resp); err != nil {
		return nil, err
	}
	var group models.Group
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupsService) Update(ctx context.Context, params models.GroupUpdateParams) (*models.Group, error) {
	var resp models.APIResponse
	if err := s.client.Post(ctx, "groups.update", params, &resp); err != nil {
		return nil, err
	}
	var group models.Group
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupsService) Delete(ctx context.Context, id string) error {
	return s.client.Post(ctx, "groups.delete", map[string]string{"id": id}, nil)
}

func (s *GroupsService) Members(ctx context.Context, id string, params models.PaginationParams) ([]models.GroupMembership, *models.Pagination, error) {
	var resp models.APIResponse
	body := map[string]interface{}{"id": id, "offset": params.Offset, "limit": params.Limit}
	if err := s.client.Post(ctx, "groups.memberships", body, &resp); err != nil {
		return nil, nil, err
	}
	var members []models.GroupMembership
	if err := json.Unmarshal(resp.Data, &members); err != nil {
		return nil, nil, err
	}
	return members, resp.Pagination, nil
}

func (s *GroupsService) AddUser(ctx context.Context, groupID, userID string) error {
	return s.client.Post(ctx, "groups.add_user", map[string]string{"id": groupID, "userId": userID}, nil)
}

func (s *GroupsService) RemoveUser(ctx context.Context, groupID, userID string) error {
	return s.client.Post(ctx, "groups.remove_user", map[string]string{"id": groupID, "userId": userID}, nil)
}
