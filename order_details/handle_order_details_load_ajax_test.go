package order_details_test

import "testing"

func TestLoadOrderDetailsSuccess(t *testing.T) {
	env := newTestEnv(t)
	order := env.seedOrder(t, "pending")

	resp := postAction(t, env, "load-order-details", nil, map[string]any{"order_id": order.GetID()})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	orderData, _ := data["order"].(map[string]any)
	if orderData["id"] != order.GetID() || orderData["status"] != "pending" {
		t.Fatalf("unexpected order data %+v", orderData)
	}
	items, _ := orderData["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 line item, got %d", len(items))
	}
	item, _ := items[0].(map[string]any)
	// line total = 9.99 * 2
	if item["total"] != 19.98 {
		t.Fatalf("expected line total 19.98, got %v", item["total"])
	}
}

func TestLoadOrderDetailsMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "load-order-details", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestLoadOrderDetailsNotFound(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "load-order-details", nil, map[string]any{"order_id": "missing"})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestLoadOrderDetailsInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "load-order-details", "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
