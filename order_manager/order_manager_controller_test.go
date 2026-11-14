package order_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOrderManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=orders", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.OrderManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Orders") {
		t.Fatal("expected page heading 'Orders'")
	}
}
