package product_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProductManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=products", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Products") {
		t.Fatal("expected page heading 'Products'")
	}
}
