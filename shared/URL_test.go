package shared

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLBuildsControllerParam(t *testing.T) {
	got := URL("/admin/shop", "products", map[string]string{"page": "2"})
	if !strings.HasPrefix(got, "/admin/shop?") {
		t.Fatalf("expected query URL, got %q", got)
	}
	if !strings.Contains(got, "controller=products") || !strings.Contains(got, "page=2") {
		t.Fatalf("missing params in %q", got)
	}
}

func TestURLDoesNotMutateParams(t *testing.T) {
	params := map[string]string{"a": "1"}
	_ = URL("/e", "home", params)
	if _, ok := params["controller"]; ok {
		t.Fatal("expected params map not to be mutated")
	}
}

func TestURLNilParams(t *testing.T) {
	got := URL("/e", "home", nil)
	if !strings.Contains(got, "controller=home") {
		t.Fatalf("expected controller param, got %q", got)
	}
}

func TestURLRUsesRequestEndpoint(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(r.Context(), KeyEndpoint, "/admin/shop")
	got := URLR(r.WithContext(ctx), "orders", nil)
	if !strings.HasPrefix(got, "/admin/shop?") || !strings.Contains(got, "controller=orders") {
		t.Fatalf("unexpected URL %q", got)
	}
}
