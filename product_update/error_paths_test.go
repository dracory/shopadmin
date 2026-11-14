package product_update_test

import (
	"context"
	"net/url"
	"testing"
)

// Error paths for the load/save actions: non-existent product and
// invalid bodies exercise writeJSONError.

func TestActionsProductNotFound(t *testing.T) {
	env := newTestEnv(t)

	for _, action := range []string{
		"load-metadata", "save-metadata",
		"load-tags", "save-tags", "load-categories", "save-categories",
		"load-status", "save-status", "load-variants", "save-variants",
		"load-notes", "save-notes", "load-details", "save-details",
	} {
		resp := postAction(t, env, action,
			url.Values{"product_id": {"missing"}}, map[string]any{})
		if resp.Status != "error" {
			t.Fatalf("action %q: expected error, got %+v", action, resp)
		}
	}
}

func TestSaveVariantsInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postActionRaw(t, env, "save-variants",
		url.Values{"product_id": {p.GetID()}}, "{oops")
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestSaveVariantsEmptyDimensionSkipped(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-variants",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{
			"dimensions": []map[string]any{
				{"name": "", "options": []string{"x"}},        // skipped: no name
				{"name": "Size", "options": []string{}},       // skipped: no options
				{"name": "Color", "options": []string{"red"}}, // kept
			},
			"variants": []map[string]any{},
		})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestSaveActionsInvalidBody(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	for _, action := range []string{
		"save-metadata", "save-tags", "save-categories",
		"save-status", "save-notes", "save-details",
	} {
		resp := postActionRaw(t, env, action,
			url.Values{"product_id": {p.GetID()}}, "{oops")
		if resp.Status != "error" {
			t.Fatalf("action %q: expected error, got %+v", action, resp)
		}
	}
}

func TestSaveVariantsUpdateAndDelete(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	dims := []map[string]any{
		{"name": "Size", "required": true, "options": []string{"S", "M"}},
	}

	// First save creates two variants.
	resp := postAction(t, env, "save-variants",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{
			"dimensions": dims,
			"variants": []map[string]any{
				{"values": map[string]string{"Size": "S"}, "sku": "W-S", "price": "9.99", "quantity": "3"},
				{"values": map[string]string{"Size": "M"}, "sku": "W-M", "price": "10.99", "quantity": "2"},
			},
		})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	variants, err := env.store.ProductVariantList(context.Background(), p.GetID())
	if err != nil || len(variants) != 2 {
		t.Fatalf("expected 2 variants, got %d (%v)", len(variants), err)
	}

	// Second save updates S, removes M.
	resp = postAction(t, env, "save-variants",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{
			"dimensions": dims,
			"variants": []map[string]any{
				{"values": map[string]string{"Size": "S"}, "sku": "W-S2", "price": "8.99", "quantity": "5"},
			},
		})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	variants, err = env.store.ProductVariantList(context.Background(), p.GetID())
	if err != nil {
		t.Fatalf("list variants: %v", err)
	}
	if len(variants) != 1 {
		t.Fatalf("expected 1 variant after removal, got %d", len(variants))
	}
	if variants[0].GetPrice() != "8.99" {
		t.Fatalf("expected updated price, got %q", variants[0].GetPrice())
	}
}
