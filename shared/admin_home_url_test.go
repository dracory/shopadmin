package shared

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestAdminHomeURL(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if got := AdminHomeURL(r); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	ctx := context.WithValue(r.Context(), KeyAdminHomeURL, "/admin")
	if got := AdminHomeURL(r.WithContext(ctx)); got != "/admin" {
		t.Fatalf("expected /admin, got %q", got)
	}
}
