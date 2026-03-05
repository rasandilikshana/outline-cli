package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client

	Documents   *DocumentsService
	Collections *CollectionsService
	Users       *UsersService
	Groups      *GroupsService
	Comments    *CommentsService
	Shares      *SharesService
	Stars       *StarsService
	Events      *EventsService
	Search      *SearchService
	Attachments *AttachmentsService
	Revisions   *RevisionsService
	Auth        *AuthService
}

func NewClient(baseURL, apiKey string) *Client {
	c := &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
	c.Documents = &DocumentsService{client: c}
	c.Collections = &CollectionsService{client: c}
	c.Users = &UsersService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Comments = &CommentsService{client: c}
	c.Shares = &SharesService{client: c}
	c.Stars = &StarsService{client: c}
	c.Events = &EventsService{client: c}
	c.Search = &SearchService{client: c}
	c.Attachments = &AttachmentsService{client: c}
	c.Revisions = &RevisionsService{client: c}
	c.Auth = &AuthService{client: c}
	return c
}

func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	} else {
		reqBody = bytes.NewReader([]byte("{}"))
	}

	reqURL := c.BaseURL + "/api/" + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Message != "" {
			return fmt.Errorf("API error %d: %s - %s", resp.StatusCode, errResp.Error, errResp.Message)
		}
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}
