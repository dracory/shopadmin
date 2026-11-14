package category_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCategoryManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=categories", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.CategoryManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Categories") {
		t.Fatal("expected page heading 'Categories'")
	}
}
