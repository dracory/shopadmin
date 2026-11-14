package discount_manager_test

import "testing"

// applyConditions is unexported; its behavior is exercised through the
// load-discounts action with conditions in the request body.

func TestLoadDiscountsConditions(t *testing.T) {
	env := newTestEnv(t)
	env.seedDiscount(t, "SUMMER10", "active")
	env.seedDiscount(t, "WINTER20", "active")
	env.seedDiscount(t, "OLDCODE", "draft")

	// code condition
	resp := postAction(t, env, "load-discounts", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "code", "operator": "=", "value": "WINTER20"},
		},
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	items, _ := resp.data()["discounts"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 discount matching code, got %d", len(items))
	}

	// status condition
	resp = postAction(t, env, "load-discounts", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "status", "operator": "=", "value": "draft"},
		},
	})
	items, _ = resp.data()["discounts"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 draft discount, got %d", len(items))
	}
}
