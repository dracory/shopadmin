package shared

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestShopAdminURL(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if got := ShopAdminURL(r); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	ctx := context.WithValue(r.Context(), KeyShopAdminURL, "/admin/shop")
	if got := ShopAdminURL(r.WithContext(ctx)); got != "/admin/shop" {
		t.Fatalf("expected /admin/shop, got %q", got)
	}
}
