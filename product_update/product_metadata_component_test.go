package product_update_test

import (
	"net/url"
	"testing"
)

// The metadata component is exercised through the controller's
// load-metadata and save-metadata actions.

func TestSaveAndLoadMetadata(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-metadata",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{"metadata": []map[string]any{
			{"key": "color", "value": "red"},
		}})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	resp = postAction(t, env, "load-metadata",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
