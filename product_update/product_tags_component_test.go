package product_update_test

import (
	"context"
	"net/url"
	"testing"
)

// The tags component is exercised through the controller's
// load-tags and save-tags actions.

func TestSaveAndLoadTags(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-tags",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{"tags": []string{"new", "sale"}})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	resp = postAction(t, env, "load-tags",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	_ = context.Background()
}
