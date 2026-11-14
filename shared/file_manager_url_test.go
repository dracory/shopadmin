package shared

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestFileManagerURL(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if got := FileManagerURL(r); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	ctx := context.WithValue(r.Context(), KeyFileManagerURL, "/admin/files")
	if got := FileManagerURL(r.WithContext(ctx)); got != "/admin/files" {
		t.Fatalf("expected /admin/files, got %q", got)
	}
}
