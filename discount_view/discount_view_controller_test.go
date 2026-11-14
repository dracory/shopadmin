package discount_view_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDiscountViewPageRenders(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "VIEW10", "active")

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=discount_view&discount_id="+d.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Discount") {
		t.Fatalf("expected discount page, got %q", rec.Body.String())
	}
}

func TestDiscountViewPageRequiresID(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=discount_view", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountView(rec, req)

	if !strings.Contains(rec.Body.String(), "required") {
		t.Fatalf("expected missing-ID alert, got %q", rec.Body.String())
	}
}
