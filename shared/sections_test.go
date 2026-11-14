package shared

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/shopstore"
)

func TestProcessProductSectionsGroupsByPosition(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	p := shopstore.NewProduct()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	res := ProcessProductSections([]ProductSectionRenderer{
		func(r *http.Request, p shopstore.ProductInterface) (ProductSection, error) {
			return ProductSection{ID: "a", Position: PositionTop, HTML: "<div>A</div>", CSS: ".a{}", JS: "a()", CSSURLs: []string{"https://x/a.css"}, JSURLs: []string{"https://x/a.js"}}, nil
		},
		func(r *http.Request, p shopstore.ProductInterface) (ProductSection, error) {
			return ProductSection{ID: "b", HTML: "<div>B</div>"}, nil // empty position -> bottom
		},
	}, r, p, logger)

	if tag := res.RenderPosition(PositionTop); tag == nil || !strings.Contains(tag.ToHTML(), "<div>A</div>") {
		t.Fatal("expected section A at top position")
	}
	if tag := res.RenderPosition(PositionBottom); tag == nil || !strings.Contains(tag.ToHTML(), "<div>B</div>") {
		t.Fatal("expected section B at bottom position")
	}
	if tag := res.RenderPosition(PositionSidebar); tag != nil {
		t.Fatal("expected no sidebar sections")
	}
	if len(res.Options.Styles) != 1 || res.Options.Styles[0] != ".a{}" {
		t.Fatalf("expected aggregated CSS, got %v", res.Options.Styles)
	}
	if len(res.Options.Scripts) != 1 || res.Options.Scripts[0] != "a()" {
		t.Fatalf("expected aggregated JS, got %v", res.Options.Scripts)
	}
	if len(res.Options.StyleURLs) != 1 || len(res.Options.ScriptURLs) != 1 {
		t.Fatalf("expected aggregated URLs, got %+v", res.Options)
	}
}

func TestProcessProductSectionsNilAndError(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	p := shopstore.NewProduct()

	res := ProcessProductSections([]ProductSectionRenderer{
		nil,
		func(r *http.Request, p shopstore.ProductInterface) (ProductSection, error) {
			return ProductSection{}, errors.New("boom")
		},
		func(r *http.Request, p shopstore.ProductInterface) (ProductSection, error) {
			return ProductSection{Position: PositionSidebar, HTML: "ok"}, nil
		},
	}, r, p, slog.New(slog.NewTextHandler(os.Stderr, nil)))

	if tag := res.RenderPosition(PositionSidebar); tag == nil || !strings.Contains(tag.ToHTML(), "ok") {
		t.Fatal("expected valid section despite nil/error renderers")
	}
}

func TestProcessProductSectionsDeduplicatesURLs(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	p := shopstore.NewProduct()

	renderer := func(r *http.Request, p shopstore.ProductInterface) (ProductSection, error) {
		return ProductSection{CSSURLs: []string{"https://x/dup.css"}, JSURLs: []string{"https://x/dup.js"}}, nil
	}
	res := ProcessProductSections([]ProductSectionRenderer{renderer, renderer}, r, p, nil)

	if len(res.Options.StyleURLs) != 1 || len(res.Options.ScriptURLs) != 1 {
		t.Fatalf("expected deduplicated URLs, got %+v", res.Options)
	}
}
