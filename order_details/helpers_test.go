package order_details_test

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

	"github.com/dracory/shopadmin/order_details"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// testEnv bundles an order details controller UI backed by an
// in-memory SQLite shopstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    order_details.UiInterface
	store shopstore.StoreInterface
}

// newTestEnv boots an in-memory SQLite shopstore and an order
// details controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	ui := order_details.UI(shared.UiConfig{
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

// seedOrder creates an order with a single line item and returns it.
func (e *testEnv) seedOrder(t *testing.T, status string) shopstore.OrderInterface {
	t.Helper()

	ctx := context.Background()

	product := shopstore.NewProduct()
	product.SetTitle("Widget")
	product.SetPrice("9.99")
	product.SetStatus("active")
	if err := e.store.ProductCreate(ctx, product); err != nil {
		t.Fatalf("seed product: %v", err)
	}

	order := shopstore.NewOrder()
	order.SetCustomerID("cust_1")
	order.SetStatus(status)
	if err := e.store.OrderCreate(ctx, order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	item := shopstore.NewOrderLineItem()
	item.SetOrderID(order.GetID())
	item.SetProductID(product.GetID())
	item.SetTitle("Widget")
	item.SetPrice("9.99")
	item.SetQuantityInt(2)
	item.SetStatus(status)
	if err := e.store.OrderLineItemCreate(ctx, item); err != nil {
		t.Fatalf("seed line item: %v", err)
	}

	return order
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
// action and JSON body, executes it, and returns the decoded envelope.
func postAction(t *testing.T, env *testEnv, action string, extraQuery url.Values, payload any) apiResp {
	t.Helper()

	q := url.Values{"action": {action}}
	for k, vs := range extraQuery {
		q[k] = vs
	}
	target := "/admin/shop?controller=order_details&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.OrderDetails(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}

// postRaw posts a raw (possibly malformed) body for the given action.
func postRaw(t *testing.T, env *testEnv, action, raw string) apiResp {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=order_details&action="+action,
		strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.OrderDetails(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
