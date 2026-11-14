package product_view_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProductViewPageRenders(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t, "Viewable", "active")

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=product_view&product_id="+p.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), p.GetID()) {
		t.Fatal("expected product ID embedded in page")
	}
}

func TestProductViewPageRequiresID(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=product_view", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductView(rec, req)

	if !strings.Contains(rec.Body.String(), "required") {
		t.Fatalf("expected missing-ID alert, got %q", rec.Body.String())
	}
}
