package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hassek/bc-cli/config"
)

func TestGetArticleUnauthenticated(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/articles/test-id/" {
			t.Errorf("Expected path /api/core/v1/content/articles/test-id/, got %s", r.URL.Path)
		}

		// Verify method
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify no auth header
		if r.Header.Get("Authorization") != "" {
			t.Errorf("Expected no Authorization header for unauthenticated request")
		}

		// Return article without bookmark status
		response := map[string]any{
			"meta": map[string]any{
				"code":    200,
				"message": "Success",
			},
			"data": map[string]any{
				"id":            "test-id",
				"category_id":   "cat-1",
				"section_id":    "sec-1",
				"title":         "Test Article",
				"summary":       "Test summary",
				"content":       "Test content",
				"author":        "Test Author",
				"read_time":     5,
				"tags":          "coffee,brewing",
				"published_at":  "2024-01-01",
				"is_bookmarked": false,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create unauthenticated client
	cfg := &config.Config{
		APIURL: server.URL,
	}
	client := NewClient(cfg)

	// Get article
	article, err := client.GetArticle("test-id")
	if err != nil {
		t.Fatalf("GetArticle failed: %v", err)
	}

	// Verify article data
	if article.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", article.ID)
	}
	if article.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got %s", article.Title)
	}
	if article.IsBookmarked {
		t.Errorf("Expected IsBookmarked false for unauthenticated user, got true")
	}
}

func TestGetArticleAuthenticated(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/articles/test-id/" {
			t.Errorf("Expected path /api/core/v1/content/articles/test-id/, got %s", r.URL.Path)
		}

		// Verify method
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify auth header is present
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("Expected Authorization header 'Bearer test-access-token', got '%s'", authHeader)
		}

		// Return article with bookmark status
		response := map[string]any{
			"meta": map[string]any{
				"code":    200,
				"message": "Success",
			},
			"data": map[string]any{
				"id":            "test-id",
				"category_id":   "cat-1",
				"section_id":    "sec-1",
				"title":         "Test Article",
				"summary":       "Test summary",
				"content":       "Test content",
				"author":        "Test Author",
				"read_time":     5,
				"tags":          "coffee,brewing",
				"published_at":  "2024-01-01",
				"is_bookmarked": true, // User has bookmarked this
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create authenticated client
	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-access-token",
	}
	client := NewClient(cfg)

	// Get article
	article, err := client.GetArticle("test-id")
	if err != nil {
		t.Fatalf("GetArticle failed: %v", err)
	}

	// Verify article data
	if article.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", article.ID)
	}
	if article.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got %s", article.Title)
	}
	if !article.IsBookmarked {
		t.Errorf("Expected IsBookmarked true for bookmarked article, got false")
	}
}

func TestListCategoryArticlesUnauthenticated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/categories/test-category/articles/" {
			t.Errorf("Expected path /api/core/v1/content/categories/test-category/articles/, got %s", r.URL.Path)
		}

		// Verify no auth header
		if r.Header.Get("Authorization") != "" {
			t.Errorf("Expected no Authorization header for unauthenticated request")
		}

		response := map[string]any{
			"meta": map[string]any{
				"code":    200,
				"message": "Success",
			},
			"data": []map[string]any{
				{
					"id":            "article-1",
					"category_id":   "test-category",
					"title":         "Article 1",
					"summary":       "Summary 1",
					"is_bookmarked": false,
				},
				{
					"id":            "article-2",
					"category_id":   "test-category",
					"title":         "Article 2",
					"summary":       "Summary 2",
					"is_bookmarked": false,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{APIURL: server.URL}
	client := NewClient(cfg)

	articles, err := client.ListCategoryArticles("test-category")
	if err != nil {
		t.Fatalf("ListCategoryArticles failed: %v", err)
	}

	if len(articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(articles))
	}

	for _, article := range articles {
		if article.IsBookmarked {
			t.Errorf("Expected IsBookmarked false for unauthenticated user, got true")
		}
	}
}

func TestListCategoryArticlesAuthenticated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header is present
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("Expected Authorization header, got '%s'", authHeader)
		}

		response := map[string]any{
			"meta": map[string]any{
				"code":    200,
				"message": "Success",
			},
			"data": []map[string]any{
				{
					"id":            "article-1",
					"category_id":   "test-category",
					"title":         "Article 1",
					"summary":       "Summary 1",
					"is_bookmarked": true, // User bookmarked this
				},
				{
					"id":            "article-2",
					"category_id":   "test-category",
					"title":         "Article 2",
					"summary":       "Summary 2",
					"is_bookmarked": false,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-access-token",
	}
	client := NewClient(cfg)

	articles, err := client.ListCategoryArticles("test-category")
	if err != nil {
		t.Fatalf("ListCategoryArticles failed: %v", err)
	}

	if len(articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(articles))
	}

	if !articles[0].IsBookmarked {
		t.Errorf("Expected first article to be bookmarked")
	}

	if articles[1].IsBookmarked {
		t.Errorf("Expected second article to not be bookmarked")
	}
}

func TestListSectionArticlesAuthenticated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/sections/sec-1/articles/" {
			t.Errorf("Expected path /api/core/v1/content/sections/sec-1/articles/, got %s", r.URL.Path)
		}

		// Verify auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("Expected Authorization header, got '%s'", authHeader)
		}

		response := map[string]any{
			"meta": map[string]any{
				"code":    200,
				"message": "Success",
			},
			"data": []map[string]any{
				{
					"id":            "article-1",
					"section_id":    "sec-1",
					"title":         "Section Article",
					"is_bookmarked": true,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-access-token",
	}
	client := NewClient(cfg)

	articles, err := client.ListSectionArticles("sec-1")
	if err != nil {
		t.Fatalf("ListSectionArticles failed: %v", err)
	}

	if len(articles) != 1 {
		t.Errorf("Expected 1 article, got %d", len(articles))
	}

	if !articles[0].IsBookmarked {
		t.Errorf("Expected article to be bookmarked")
	}
}

func TestCreateBookmark(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/bookmarks/" {
			t.Errorf("Expected path /api/core/v1/content/bookmarks/, got %s", r.URL.Path)
		}

		// Verify method
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Verify auth required
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("Expected Authorization header, got '%s'", authHeader)
		}

		// Parse request body
		var requestBody map[string]string
		_ = json.NewDecoder(r.Body).Decode(&requestBody)

		if requestBody["article_id"] != "article-123" {
			t.Errorf("Expected article_id 'article-123', got '%s'", requestBody["article_id"])
		}

		response := map[string]any{
			"meta": map[string]any{
				"code":    201,
				"message": "Created",
			},
			"data": map[string]any{
				"id":         "bookmark-1",
				"article_id": "article-123",
				"created_at": "2024-01-01T00:00:00Z",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-access-token",
	}
	client := NewClient(cfg)

	bookmark, err := client.CreateBookmark("article-123")
	if err != nil {
		t.Fatalf("CreateBookmark failed: %v", err)
	}

	if bookmark.ID != "bookmark-1" {
		t.Errorf("Expected bookmark ID 'bookmark-1', got '%s'", bookmark.ID)
	}

	if bookmark.ArticleID != "article-123" {
		t.Errorf("Expected article ID 'article-123', got '%s'", bookmark.ArticleID)
	}
}

func TestDeleteBookmark(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/bookmarks/bookmark-1/" {
			t.Errorf("Expected path /api/core/v1/content/bookmarks/bookmark-1/, got %s", r.URL.Path)
		}

		// Verify method
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}

		// Verify auth required
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("Expected Authorization header, got '%s'", authHeader)
		}

		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-access-token",
	}
	client := NewClient(cfg)

	err := client.DeleteBookmark("bookmark-1")
	if err != nil {
		t.Fatalf("DeleteBookmark failed: %v", err)
	}
}

func TestListBookmarks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/api/core/v1/content/bookmarks/" {
			t.Errorf("Expected path /api/core/v1/content/bookmarks/, got %s", r.URL.Path)
		}

		// Verify method
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify auth required
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("Expected Authorization header, got '%s'", authHeader)
		}

		response := map[string]any{
			"meta": map[string]any{
				"code":    200,
				"message": "Success",
			},
			"data": map[string]any{
				"count": 2,
				"results": []map[string]any{
					{
						"id": "bookmark-1",
						// Note: article_id not populated by backend, use nested article.id instead
						"article": map[string]any{
							"id":            "article-1",
							"title":         "Bookmarked Article 1",
							"is_bookmarked": true,
						},
					},
					{
						"id": "bookmark-2",
						"article": map[string]any{
							"id":            "article-2",
							"title":         "Bookmarked Article 2",
							"is_bookmarked": true,
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-access-token",
	}
	client := NewClient(cfg)

	bookmarks, err := client.ListBookmarks()
	if err != nil {
		t.Fatalf("ListBookmarks failed: %v", err)
	}

	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}

	if bookmarks[0].ID != "bookmark-1" {
		t.Errorf("Expected bookmark ID 'bookmark-1', got '%s'", bookmarks[0].ID)
	}

	if bookmarks[0].Article.ID != "article-1" {
		t.Errorf("Expected article ID 'article-1', got '%s'", bookmarks[0].Article.ID)
	}

	if !bookmarks[0].Article.IsBookmarked {
		t.Errorf("Expected bookmarked article to have IsBookmarked true")
	}
}
