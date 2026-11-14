package shared

import (
	"strings"
	"testing"
)

func TestErrorAlert(t *testing.T) {
	html := ErrorAlert("something broke")
	if !strings.Contains(html, "alert alert-danger") {
		t.Fatalf("expected danger alert class, got %q", html)
	}
	if !strings.Contains(html, "something broke") {
		t.Fatalf("expected message in output, got %q", html)
	}
}
