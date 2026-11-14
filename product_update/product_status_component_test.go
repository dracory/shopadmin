package product_update_test

import (
	"context"
	"net/url"
	"testing"
)

// The status component is exercised through the controller's
// load-status and save-status actions.

func TestLoadStatus(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "load-status",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestSaveStatus(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-status",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{"status": "draft"})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	updated, _ := env.store.ProductFindByID(context.Background(), p.GetID())
	if updated.GetStatus() != "draft" {
		t.Fatalf("expected status draft, got %q", updated.GetStatus())
	}
}
