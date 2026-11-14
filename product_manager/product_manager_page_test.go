package product_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProductManagerPageContentType(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=products", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductManager(rec, req)

	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}
}
