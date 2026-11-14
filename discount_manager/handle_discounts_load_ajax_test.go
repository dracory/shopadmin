package discount_manager_test

import "testing"

func TestLoadDiscountsSuccess(t *testing.T) {
	env := newTestEnv(t)
	env.seedDiscount(t, "C1", "active")
	env.seedDiscount(t, "C2", "draft")

	resp := postAction(t, env, "load-discounts", nil, map[string]any{})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if total, _ := data["total"].(float64); total != 2 {
		t.Fatalf("expected total 2, got %v", data["total"])
	}
	items, _ := data["discounts"].([]any)
	if len(items) != 2 {
		t.Fatalf("expected 2 discounts, got %d", len(items))
	}
}

func TestLoadDiscountsStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	env.seedDiscount(t, "C1", "active")
	env.seedDiscount(t, "C2", "draft")

	resp := postAction(t, env, "load-discounts", nil, map[string]any{"status": "draft"})
	items, _ := resp.data()["discounts"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 draft discount, got %d", len(items))
	}
}

func TestLoadDiscountsInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "load-discounts", "{bad")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
