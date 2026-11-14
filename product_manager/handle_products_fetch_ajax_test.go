package product_manager_test

import "testing"

func TestLoadProductsSuccess(t *testing.T) {
	env := newTestEnv(t)
	env.seedProduct(t, "Alpha", "active")
	env.seedProduct(t, "Beta", "draft")

	resp := postAction(t, env, "load-products", nil, map[string]any{})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if total, _ := data["total"].(float64); total != 2 {
		t.Fatalf("expected total 2, got %v", data["total"])
	}
	items, _ := data["products"].([]any)
	if len(items) != 2 {
		t.Fatalf("expected 2 products, got %d", len(items))
	}
}

func TestLoadProductsStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	env.seedProduct(t, "Alpha", "active")
	env.seedProduct(t, "Beta", "draft")

	resp := postAction(t, env, "load-products", nil, map[string]any{"status": "draft"})
	items, _ := resp.data()["products"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 draft product, got %d", len(items))
	}
}

func TestLoadProductsInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "load-products", "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
