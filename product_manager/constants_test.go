package product_manager_test

import "testing"

// applyConditions / categoryProductIDs are unexported; they are
// exercised through the load-products action with conditions.

func TestLoadProductsConditions(t *testing.T) {
	env := newTestEnv(t)
	env.seedProduct(t, "Hammer", "active")
	env.seedProduct(t, "Drill", "active")
	env.seedProduct(t, "Old Drill", "draft")

	// search condition filters by title
	resp := postAction(t, env, "load-products", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "search", "operator": "=", "value": "Drill"},
		},
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	items, _ := resp.data()["products"].([]any)
	if len(items) != 2 {
		t.Fatalf("expected 2 products matching 'Drill', got %d", len(items))
	}

	// status condition
	resp = postAction(t, env, "load-products", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "status", "operator": "=", "value": "draft"},
		},
	})
	items, _ = resp.data()["products"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 draft product, got %d", len(items))
	}
}
