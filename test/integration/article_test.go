package integration

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/hexlet-components/go-gin-example/db/generated"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "create article success",
			body:           `{"name":"Test Article"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"name":"Test Article"}`,
		},
		{
			name:           "create article empty body",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "create article invalid json",
			body:           `{"name":}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "create article missing name",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter(t)

			var body *bytes.Buffer
			if tt.body != "" {
				body = bytes.NewBufferString(tt.body)
			} else {
				body = bytes.NewBuffer(nil)
			}

			req, _ := http.NewRequest("POST", "/articles", body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetArticle(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
		setup          func(queries *db.Queries) int64
	}{
		{
			name:           "get article success",
			url:            "/articles/1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Test Article"}`,
			setup: func(queries *db.Queries) int64 {
				article, err := queries.CreateArticle(context.Background(), "Test Article")
				if err != nil {
					t.Fatalf("failed to create test article: %v", err)
				}
				return article.ID
			},
		},
		{
			name:           "get article not found",
			url:            "/articles/999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "get article invalid id",
			url:            "/articles/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "get article negative id",
			url:            "/articles/-1",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, queries := setupTestRouterWithQueries(t)

			if tt.setup != nil {
				tt.setup(queries)
			}

			req, _ := http.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestListArticles(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
		setup          func(*db.Queries)
	}{
		{
			name:           "list articles empty",
			expectedStatus: http.StatusOK,
			expectedBody:   "null",
		},
		{
			name:           "list articles with single article",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"Single Article"}]`,
			setup: func(queries *db.Queries) {
				_, err := queries.CreateArticle(context.Background(), "Single Article")
				if err != nil {
					t.Fatalf("failed to create test article: %v", err)
				}
			},
		},
		{
			name:           "list articles with multiple articles",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"Article 1"},{"id":2,"name":"Article 2"}]`,
			setup: func(queries *db.Queries) {
				for _, name := range []string{"Article 1", "Article 2"} {
					_, err := queries.CreateArticle(context.Background(), name)
					if err != nil {
						t.Fatalf("failed to create test article: %v", err)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, queries := setupTestRouterWithQueries(t)

			if tt.setup != nil {
				tt.setup(queries)
			}

			req, _ := http.NewRequest("GET", "/articles", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		body           string
		expectedStatus int
		expectedBody   string
		setup          func(*db.Queries) int64
	}{
		{
			name:           "update article success",
			url:            "/articles/1",
			body:           `{"name":"Updated Article"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Updated Article"}`,
			setup: func(queries *db.Queries) int64 {
				article, err := queries.CreateArticle(context.Background(), "Original Article")
				if err != nil {
					t.Fatalf("failed to create test article: %v", err)
				}
				return article.ID
			},
		},
		{
			name:           "update article not found",
			url:            "/articles/999",
			body:           `{"name":"Updated Article"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "update article invalid id",
			url:            "/articles/abc",
			body:           `{"name":"Updated Article"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "update article empty body",
			url:            "/articles/1",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			setup: func(queries *db.Queries) int64 {
				article, err := queries.CreateArticle(context.Background(), "Original Article")
				if err != nil {
					t.Fatalf("failed to create test article: %v", err)
				}
				return article.ID
			},
		},
		{
			name:           "update article invalid json",
			url:            "/articles/1",
			body:           `{"name":}`,
			expectedStatus: http.StatusBadRequest,
			setup: func(queries *db.Queries) int64 {
				article, err := queries.CreateArticle(context.Background(), "Original Article")
				if err != nil {
					t.Fatalf("failed to create test article: %v", err)
				}
				return article.ID
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, queries := setupTestRouterWithQueries(t)

			if tt.setup != nil {
				tt.setup(queries)
			}

			var body *bytes.Buffer
			if tt.body != "" {
				body = bytes.NewBufferString(tt.body)
			} else {
				body = bytes.NewBuffer(nil)
			}

			req, _ := http.NewRequest("PUT", tt.url, body)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

// TestDeleteArticle tests for Delete operations
func TestDeleteArticle(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		setup          func(*db.Queries) int64
	}{
		{
			name:           "delete article success",
			url:            "/articles/1",
			expectedStatus: http.StatusNoContent,
			setup: func(queries *db.Queries) int64 {
				article, err := queries.CreateArticle(context.Background(), "To Delete")
				if err != nil {
					t.Fatalf("failed to create test article: %v", err)
				}
				return article.ID
			},
		},
		{
			name:           "delete article not found",
			url:            "/articles/999",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete article invalid id",
			url:            "/articles/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "delete article negative id",
			url:            "/articles/-1",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, queries := setupTestRouterWithQueries(t)

			if tt.setup != nil {
				tt.setup(queries)
			}

			req, _ := http.NewRequest("DELETE", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
