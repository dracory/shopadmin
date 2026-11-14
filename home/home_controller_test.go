package home_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/shopadmin/home"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// testEnv bundles a home controller UI backed by an in-memory SQLite
// shopstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    home.UiInterface
	store shopstore.StoreInterface
}

// newTestEnv boots an in-memory SQLite shopstore and a home
// controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	ui := home.UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	})

	return &testEnv{ui: ui, store: store}
}

func TestHomePageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=home", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.Home(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Shop. Dashboard") {
		t.Fatalf("expected dashboard heading, got %q", rec.Body.String())
	}
}

func TestLoadStats(t *testing.T) {
	env := newTestEnv(t)

	p := shopstore.NewProduct()
	p.SetTitle("P1")
	p.SetStatus("active")
	if err := env.store.ProductCreate(context.Background(), p); err != nil {
		t.Fatalf("seed product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=home&action=load-stats", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.Home(rec, req)

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	if resp["status"] != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data, _ := resp["data"].(map[string]any)
	if data["product_count"] != float64(1) {
		t.Fatalf("expected product_count 1, got %v", data["product_count"])
	}
	for _, key := range []string{"category_count", "discount_count", "order_count"} {
		if _, ok := data[key]; !ok {
			t.Fatalf("expected %q in stats", key)
		}
	}
}
