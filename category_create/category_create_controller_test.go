package category_create_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCategoryCreatePageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=category_create", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.CategoryCreate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"Create Category", "category-title", "category-status", "btn-create-category"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected %q in page body", want)
		}
	}
}

func TestCreateCategorySuccess(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "create-category", nil, map[string]any{
		"title":       "New Cat",
		"description": "desc",
		"status":      "active",
		"parent_id":   "",
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	id, _ := resp.data()["category_id"].(string)
	if id == "" {
		t.Fatal("expected category_id in response data")
	}

	category, err := env.store.CategoryFindByID(context.Background(), id)
	if err != nil || category == nil {
		t.Fatalf("expected category persisted, got %v", err)
	}
	if category.GetTitle() != "New Cat" || category.GetStatus() != "active" {
		t.Fatalf("unexpected category fields: %q %q", category.GetTitle(), category.GetStatus())
	}
}

func TestCreateCategoryMissingTitle(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "create-category", nil, map[string]any{"title": ""})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
	if !strings.Contains(resp.message(), "Title") {
		t.Fatalf("expected title error message, got %q", resp.message())
	}
}

func TestCreateCategoryInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodPost, "/admin/shop?controller=category_create&action=create-category", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.ui.CategoryCreate(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
