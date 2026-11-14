package order_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOrderManagerPageContentType(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=orders", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.OrderManager(rec, req)

	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %q", ct)
	}
}
