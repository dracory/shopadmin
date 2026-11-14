package shared

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLayoutRendersWebpage(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	html := Layout(w, r, "My Title", "<p>Hello</p>", struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		Styles:    []string{".x{color:red}"},
		StyleURLs: []string{"https://example.com/x.css"},
		Scripts:   []string{"console.log('y')"},
	})

	for _, want := range []string{"My Title", "<p>Hello</p>", ".x{color:red}", "https://example.com/x.css", "console.log('y')", "bootstrap"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in layout output", want)
		}
	}
}
