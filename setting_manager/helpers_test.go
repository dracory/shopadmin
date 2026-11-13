package setting_manager_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/dracory/shopadmin/setting_manager"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// testEnv bundles a setting manager controller UI backed by an in-memory
// SQLite shopstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    setting_manager.UiInterface
	store shopstore.StoreInterface
}

// newTestEnv boots an in-memory SQLite shopstore and a setting manager
// controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ui := setting_manager.UI(shared.UiConfig{
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
	target := "/admin/shop/settings?controller=settings&" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	env.ui.SettingManager(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v (body=%q)", err, rec.Body.String())
	}
	return resp
}
