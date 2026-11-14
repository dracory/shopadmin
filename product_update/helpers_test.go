package product_update_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dracory/shopadmin/product_update"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// testEnv bundles a product update controller UI backed by an in-memory
// SQLite shopstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    product_update.UiInterface
	store shopstore.StoreInterface
}

// newTestEnv boots an in-memory SQLite shopstore and a product update
// controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ui := product_update.UI(shared.UiConfig{
		Store:  store,
		Logger: logger,
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	}, "")

	return &testEnv{ui: ui, store: store}
}

// seedProduct creates a product with the common fields populated and
// returns it.
func (e *testEnv) seedProduct(t *testing.T) shopstore.ProductInterface {
	t.Helper()

	product := shopstore.NewProduct().
		SetTitle("Test Product").
		SetStatus(shopstore.PRODUCT_STATUS_ACTIVE).
		SetPrice("9.99").
		SetQuantity("5").
		SetDescription("old description").
		SetMemo("old memo")

	if err := e.store.ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("seed product: %v", err)
	}
	return product
}

type apiResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// postAction builds a POST request to the controller with the given
// action and JSON body, executes it, and returns the decoded envelope.
func postAction(t *testing.T, env *testEnv, action string, extraQuery url.Values, payload any) apiResponse {
	t.Helper()

	q := url.Values{"action": {action}}
	for k, vs := range extraQuery {
		q[k] = vs
	}
	target := "/admin/shop/products?controller=product-update&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}

// postActionRaw posts a raw (possibly malformed) JSON body for the
// given action.
func postActionRaw(t *testing.T, env *testEnv, action string, extraQuery url.Values, raw string) apiResponse {
	t.Helper()

	q := url.Values{"action": {action}}
	for k, vs := range extraQuery {
		q[k] = vs
	}

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop/products?controller=product-update&"+q.Encode(),
		strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
