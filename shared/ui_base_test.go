package shared

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewUiBaseMapsConfig(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	layout := func(w http.ResponseWriter, r *http.Request, title, html string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string {
		return "LAYOUT:" + title + ":" + html
	}

	base := NewUiBase(UiConfig{
		Logger: logger,
		Layout: layout,
	})

	if base.Logger() != logger {
		t.Fatal("expected logger to be set")
	}
	if base.Store() != nil {
		t.Fatal("expected nil store")
	}
	if base.CustomerResolver() != nil {
		t.Fatal("expected nil customer resolver")
	}
	if base.ProductViewSections() != nil || base.ProductUpdateSections() != nil {
		t.Fatal("expected nil section renderers")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	if got := base.Layout(w, r, "T", "B", struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{}); got != "LAYOUT:T:B" {
		t.Fatalf("unexpected layout output %q", got)
	}
}
