package product_update_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestProductUpdatePageRenders(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=product_update&product_id="+p.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), p.GetID()) {
		t.Fatal("expected product ID embedded in page")
	}
}

func TestProductUpdatePageRequiresID(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=product_update", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	if !strings.Contains(rec.Body.String(), "Product ID is required") {
		t.Fatalf("expected missing-ID alert, got %q", rec.Body.String())
	}
}

func TestProductUpdatePageNotFound(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=product_update&product_id=missing", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	if !strings.Contains(rec.Body.String(), "Product not found") {
		t.Fatalf("expected not-found alert, got %q", rec.Body.String())
	}
}

func TestLoadDetails(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "load-details",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	if resp.Data["title"] != "Test Product" {
		t.Fatalf("unexpected loaded data %+v", resp.Data)
	}
}
