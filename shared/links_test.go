package shared

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewLinksDefaultsBaseURL(t *testing.T) {
	l := NewLinks("")
	if got := l.Home(nil); !strings.HasPrefix(got, "/admin/shop?") {
		t.Fatalf("expected default base URL, got %q", got)
	}
}

func TestLinksBuildControllerURLs(t *testing.T) {
	l := NewLinks("/x")
	cases := map[string]string{
		"home":            l.Home(nil),
		"products":        l.Products(nil),
		"product_view":    l.ProductView(nil),
		"product_update":  l.ProductUpdate(nil),
		"categories":      l.Categories(nil),
		"category_create": l.CategoryCreate(nil),
		"category_update": l.CategoryUpdate(nil),
		"discounts":       l.Discounts(nil),
		"discount_view":   l.DiscountView(nil),
		"discount_update": l.DiscountUpdate(nil),
		"orders":          l.Orders(nil),
		"order_details":   l.OrderDetails(nil),
		"settings":        l.Settings(nil),
	}
	for controller, url := range cases {
		if !strings.HasPrefix(url, "/x?") || !strings.Contains(url, "controller="+controller) {
			t.Fatalf("unexpected URL %q for controller %q", url, controller)
		}
	}
}

func TestNewLinksFromRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(r.Context(), KeyShopAdminURL, "/custom")
	l := NewLinksFromRequest(r.WithContext(ctx))
	if got := l.Orders(nil); !strings.HasPrefix(got, "/custom?") {
		t.Fatalf("expected request base URL, got %q", got)
	}
}
