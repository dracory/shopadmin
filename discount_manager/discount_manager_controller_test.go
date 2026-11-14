package discount_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDiscountManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=discounts", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Discounts") {
		t.Fatal("expected page heading 'Discounts'")
	}
}
