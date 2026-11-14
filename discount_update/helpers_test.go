package discount_update_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/dracory/shopadmin/discount_update"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// testEnv bundles a discount update controller UI backed by an
// in-memory SQLite shopstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    discount_update.UiInterface
	store shopstore.StoreInterface
}

// newTestEnv boots an in-memory SQLite shopstore and a discount
// update controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ui := discount_update.UI(shared.UiConfig{
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
	})

	return &testEnv{ui: ui, store: store}
}

// seedDiscount creates a discount with the given code/status.
func (e *testEnv) seedDiscount(t *testing.T, code, status string) shopstore.DiscountInterface {
	t.Helper()

	discount := shopstore.NewDiscount()
	discount.SetCode(code)
	discount.SetTitle("Discount " + code)
	discount.SetStatus(status)
	if err := e.store.DiscountCreate(context.Background(), discount); err != nil {
		t.Fatalf("seed discount: %v", err)
	}
	return discount
}

// apiResp is the decoded JSON envelope returned by the api package.
type apiResp map[string]any

func (a apiResp) status() string  { s, _ := a["status"].(string); return s }
func (a apiResp) message() string { s, _ := a["message"].(string); return s }
func (a apiResp) data() map[string]any {
	d, _ := a["data"].(map[string]any)
	return d
}

// postAction builds a POST request to the controller with the given
// action, discount ID, and JSON body, and returns the decoded envelope.
func postAction(t *testing.T, env *testEnv, action, discountID string, payload any) apiResp {
	t.Helper()

	q := url.Values{"action": {action}, "discount_id": {discountID}}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=discount_update&"+q.Encode(),
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.DiscountUpdate(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
