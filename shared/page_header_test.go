package shared

import (
	"strings"
	"testing"

	"github.com/dracory/hb"
)

func TestPageHeaderTitle(t *testing.T) {
	html := PageHeader("Shop. Products").ToHTML()
	if !strings.Contains(html, "Shop. Products") || !strings.Contains(html, "<h2") {
		t.Fatalf("expected H2 title, got %q", html)
	}
}

func TestPageHeaderActions(t *testing.T) {
	btn := hb.Button().Text("Add")
	html := PageHeader("Title", btn).ToHTML()
	if !strings.Contains(html, "Add") || !strings.Contains(html, "d-flex gap-2") {
		t.Fatalf("expected action button container, got %q", html)
	}
}
