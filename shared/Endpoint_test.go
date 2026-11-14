package shared

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestEndpointFromContext(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if got := Endpoint(r); got != "" {
		t.Fatalf("expected empty endpoint, got %q", got)
	}

	ctx := context.WithValue(r.Context(), KeyEndpoint, "/admin/shop")
	if got := Endpoint(r.WithContext(ctx)); got != "/admin/shop" {
		t.Fatalf("expected /admin/shop, got %q", got)
	}
}
