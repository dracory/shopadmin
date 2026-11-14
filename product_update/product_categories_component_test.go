package product_update_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/dracory/shopstore"
)

// The categories component is exercised through the controller's
// load-categories and save-categories actions.

func TestSaveAndLoadCategories(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	cat := shopstore.NewCategory()
	cat.SetTitle("Gadgets")
	cat.SetStatus("active")
	if err := env.store.CategoryCreate(context.Background(), cat); err != nil {
		t.Fatalf("seed category: %v", err)
	}

	resp := postAction(t, env, "save-categories",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{"category_ids": []string{cat.GetID()}})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	resp = postAction(t, env, "load-categories",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
