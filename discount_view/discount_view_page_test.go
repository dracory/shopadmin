package discount_view_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// escapeJSString is unexported; the page render path is covered by
// TestDiscountViewPageRenders. This test verifies the page embeds the
// discount ID so the frontend can fetch details.
func TestDiscountViewPageEmbedsDiscountID(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "EMBED", "active")

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=discount_view&discount_id="+d.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.DiscountView(rec, req)

	if !strings.Contains(rec.Body.String(), d.GetID()) {
		t.Fatal("expected discount ID embedded in page")
	}
}
