package discount_update_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDiscountUpdatePageRenders(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "SAVE10", "active")

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=discount_update&discount_id="+d.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountUpdate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Edit Discount") {
		t.Fatalf("expected 'Edit Discount' heading, got %q", rec.Body.String())
	}
}

func TestDiscountUpdatePageRequiresID(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=discount_update", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountUpdate(rec, req)

	if !strings.Contains(rec.Body.String(), "Discount ID is required") {
		t.Fatalf("expected missing-ID alert, got %q", rec.Body.String())
	}
}

func TestDiscountUpdatePageNotFound(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=discount_update&discount_id=missing", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountUpdate(rec, req)

	if !strings.Contains(rec.Body.String(), "Discount not found") {
		t.Fatalf("expected not-found alert, got %q", rec.Body.String())
	}
}
