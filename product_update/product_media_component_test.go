package product_update_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/dracory/shopstore"
	"net/url"
	"testing"
)

// The media component is exercised through the controller's
// load-media action (upload requires multipart form data).

func TestLoadMediaEmpty(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "load-media",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestSaveMedia(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-media",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{"media": []map[string]any{
			{"url": "https://example.com/a.png", "fileName": "a.png", "type": "image/png"},
			{"url": "", "fileName": "skip.png", "type": "image/png"},
		}})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	medias, err := env.store.MediaList(context.Background(),
		shopstore.NewMediaQuery().SetEntityID(p.GetID()))
	if err != nil {
		t.Fatalf("list media: %v", err)
	}
	if len(medias) != 1 {
		t.Fatalf("expected 1 media (empty URL skipped), got %d", len(medias))
	}
}

func TestSaveMediaInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=product_update&action=save-media&product_id="+p.GetID(),
		strings.NewReader("{oops"))
	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
