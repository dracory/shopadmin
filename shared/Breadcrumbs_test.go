package shared

import (
	"strings"
	"testing"
)

func TestBreadcrumbsRendersItems(t *testing.T) {
	html := Breadcrumbs([]Breadcrumb{
		{Name: "Home", URL: "/admin"},
		{Name: "Shop", URL: "/admin/shop"},
	}).ToHTML()

	if !strings.Contains(html, "breadcrumb") {
		t.Fatalf("expected breadcrumb class, got %q", html)
	}
	if !strings.Contains(html, "Home") || !strings.Contains(html, "Shop") {
		t.Fatalf("expected item names, got %q", html)
	}
	if !strings.Contains(html, "/admin/shop") {
		t.Fatalf("expected item links, got %q", html)
	}
}

func TestBreadcrumbsEmpty(t *testing.T) {
	html := Breadcrumbs(nil).ToHTML()
	if !strings.Contains(html, "breadcrumb") {
		t.Fatalf("expected empty breadcrumb nav, got %q", html)
	}
}
