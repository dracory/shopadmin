package category_update_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCategoryUpdatePageRenders(t *testing.T) {
	env := newTestEnv(t)
	cat := env.seedCategory(t, "Editable", "active")

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=category_update&category_id="+cat.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.CategoryUpdate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"Update Category", "category-title", "btn-update-category", cat.GetID()} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected %q in page body", want)
		}
	}
}

func TestCategoryUpdatePageRequiresID(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=category_update", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.CategoryUpdate(rec, req)

	if !strings.Contains(rec.Body.String(), "Category ID is required") {
		t.Fatalf("expected missing-ID alert, got %q", rec.Body.String())
	}
}

func TestLoadCategory(t *testing.T) {
	env := newTestEnv(t)
	cat := env.seedCategory(t, "Editable", "active")

	resp := postAction(t, env, "load-category",
		url.Values{"category_id": {cat.GetID()}}, nil)
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	loaded, _ := resp.data()["category"].(map[string]any)
	if loaded["title"] != "Editable" {
		t.Fatalf("unexpected loaded category %+v", loaded)
	}
}

func TestLoadCategoryNotFound(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "load-category",
		url.Values{"category_id": {"missing"}}, nil)
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestUpdateCategorySuccess(t *testing.T) {
	env := newTestEnv(t)
	cat := env.seedCategory(t, "Old", "draft")

	resp := postAction(t, env, "update-category", nil, map[string]any{
		"category_id": cat.GetID(),
		"title":       "New Title",
		"description": "new desc",
		"status":      "active",
		"parent_id":   "parent-1",
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	updated, err := env.store.CategoryFindByID(context.Background(), cat.GetID())
	if err != nil || updated == nil {
		t.Fatalf("find updated: %v", err)
	}
	if updated.GetTitle() != "New Title" || updated.GetStatus() != "active" || updated.GetParentID() != "parent-1" {
		t.Fatalf("unexpected updated fields: %+v", updated)
	}
}

func TestUpdateCategoryMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "update-category", nil, map[string]any{"title": "x"})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestUpdateCategoryInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "update-category", nil, "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
