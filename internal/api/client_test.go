package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://example.com", "ol_api_test123")

	if client.BaseURL != "https://example.com" {
		t.Errorf("BaseURL = %q, want %q", client.BaseURL, "https://example.com")
	}
	if client.APIKey != "ol_api_test123" {
		t.Errorf("APIKey = %q, want %q", client.APIKey, "ol_api_test123")
	}
	if client.Documents == nil {
		t.Error("Documents service should not be nil")
	}
	if client.Collections == nil {
		t.Error("Collections service should not be nil")
	}
	if client.Users == nil {
		t.Error("Users service should not be nil")
	}
	if client.Groups == nil {
		t.Error("Groups service should not be nil")
	}
	if client.Comments == nil {
		t.Error("Comments service should not be nil")
	}
	if client.Shares == nil {
		t.Error("Shares service should not be nil")
	}
	if client.Stars == nil {
		t.Error("Stars service should not be nil")
	}
	if client.Events == nil {
		t.Error("Events service should not be nil")
	}
	if client.Search == nil {
		t.Error("Search service should not be nil")
	}
	if client.Attachments == nil {
		t.Error("Attachments service should not be nil")
	}
	if client.Revisions == nil {
		t.Error("Revisions service should not be nil")
	}
	if client.Auth == nil {
		t.Error("Auth service should not be nil")
	}
}

func TestNewClient_TrailingSlash(t *testing.T) {
	client := NewClient("https://example.com/", "key")
	if client.BaseURL != "https://example.com" {
		t.Errorf("BaseURL should strip trailing slash, got %q", client.BaseURL)
	}
}

func TestClient_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Authorization = %q, want %q", r.Header.Get("Authorization"), "Bearer test-key")
		}
		if r.URL.Path != "/api/test.endpoint" {
			t.Errorf("Path = %q, want /api/test.endpoint", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":   true,
			"data": map[string]string{"id": "123", "name": "test"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	var result struct {
		OK   bool `json:"ok"`
		Data struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}

	err := client.Post(context.Background(), "test.endpoint", map[string]string{"key": "value"}, &result)
	if err != nil {
		t.Fatalf("Post() error: %v", err)
	}
	if !result.OK {
		t.Error("Expected ok=true")
	}
	if result.Data.ID != "123" {
		t.Errorf("Data.ID = %q, want %q", result.Data.ID, "123")
	}
}

func TestClient_Post_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      false,
			"error":   "authentication_required",
			"message": "Authentication is required",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-key")
	err := client.Post(context.Background(), "test.endpoint", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 401 response")
	}
}

func TestClient_Post_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "key")
	err := client.Post(context.Background(), "test.endpoint", nil, nil)
	if err == nil {
		t.Fatal("Expected error for 500 response")
	}
}

func TestClient_Post_NilBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "data": nil})
	}))
	defer server.Close()

	client := NewClient(server.URL, "key")
	err := client.Post(context.Background(), "test.endpoint", nil, nil)
	if err != nil {
		t.Fatalf("Post() with nil body should succeed, got: %v", err)
	}
}

func TestClient_Post_ConnectionRefused(t *testing.T) {
	client := NewClient("http://localhost:1", "key")
	err := client.Post(context.Background(), "test.endpoint", nil, nil)
	if err == nil {
		t.Fatal("Expected error for connection refused")
	}
}
