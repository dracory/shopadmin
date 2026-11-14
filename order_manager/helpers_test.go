package order_manager_test

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

	"github.com/dracory/shopadmin/order_manager"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// testEnv bundles an order manager controller UI backed by an
// in-memory SQLite shopstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    order_manager.UiInterface
	store shopstore.StoreInterface
}

// mockCustomerResolver implements shared.CustomerResolverInterface.
type mockCustomerResolver struct {
	name  string
	email string
	ids   []string
	err   error
}

func (m *mockCustomerResolver) FindByID(_ context.Context, id string) (string, string) {
	return m.name, m.email
}

func (m *mockCustomerResolver) SearchIDs(_ context.Context, name, email string) ([]string, error) {
	return m.ids, m.err
}

// newTestEnv boots an in-memory SQLite shopstore and an order
// manager controller UI.
func newTestEnv(t *testing.T) *testEnv {
	return newTestEnvWithResolver(t, nil)
}

func newTestEnvWithResolver(t *testing.T, resolver shared.CustomerResolverInterface) *testEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	ui := order_manager.UI(shared.UiConfig{
		Store:            store,
		Logger:           slog.New(slog.NewTextHandler(os.Stderr, nil)),
		CustomerResolver: resolver,
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

// seedOrder creates an order with the given customer/status.
func (e *testEnv) seedOrder(t *testing.T, customerID, status string) shopstore.OrderInterface {
	t.Helper()

	order := shopstore.NewOrder()
	order.SetCustomerID(customerID)
	order.SetStatus(status)
	if err := e.store.OrderCreate(context.Background(), order); err != nil {
		t.Fatalf("seed order: %v", err)
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
	target := "/admin/shop?controller=orders&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.OrderManager(rec, req)

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
		"/admin/shop?controller=orders&action="+action,
		strings.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.OrderManager(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
