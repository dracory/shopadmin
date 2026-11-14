package product_update_test

import (
	"context"
	"net/url"
	"testing"
)

// The admin notes component is exercised through the controller's
// load-notes and save-notes actions.

func TestSaveAndLoadNotes(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-notes",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{"memo": "internal note"})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	updated, _ := env.store.ProductFindByID(context.Background(), p.GetID())
	if updated.GetMemo() != "internal note" {
		t.Fatalf("expected memo saved, got %q", updated.GetMemo())
	}

	resp = postAction(t, env, "load-notes",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
