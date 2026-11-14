package order_details_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOrderDetailsPageRenders(t *testing.T) {
	env := newTestEnv(t)
	order := env.seedOrder(t, "pending")

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=order_details&order_id="+order.GetID(), http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.OrderDetails(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), order.GetID()) {
		t.Fatal("expected order ID embedded in page")
	}
}

func TestOrderDetailsPageRequiresID(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=order_details", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.OrderDetails(rec, req)

	if !strings.Contains(rec.Body.String(), "required") {
		t.Fatalf("expected missing-ID alert, got %q", rec.Body.String())
	}
}
