package order_manager_test

import "testing"

func TestLoadOrdersSuccess(t *testing.T) {
	env := newTestEnv(t)
	env.seedOrder(t, "cust_1", "pending")
	env.seedOrder(t, "cust_2", "completed")

	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if total, _ := data["total"].(float64); total != 2 {
		t.Fatalf("expected total 2, got %v", data["total"])
	}
	orders, _ := data["orders"].([]any)
	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}
}

func TestLoadOrdersStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	env.seedOrder(t, "cust_1", "pending")
	env.seedOrder(t, "cust_2", "completed")

	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{"status": "pending"})
	orders, _ := resp.data()["orders"].([]any)
	if len(orders) != 1 {
		t.Fatalf("expected 1 pending order, got %d", len(orders))
	}
}

func TestLoadOrdersOrderIDFilter(t *testing.T) {
	env := newTestEnv(t)
	o := env.seedOrder(t, "cust_1", "pending")
	env.seedOrder(t, "cust_2", "completed")

	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{
		"order_id": o.GetID(),
	})
	orders, _ := resp.data()["orders"].([]any)
	if len(orders) != 1 {
		t.Fatalf("expected 1 order matching ID, got %d", len(orders))
	}
}

func TestLoadOrdersPagination(t *testing.T) {
	env := newTestEnv(t)
	for i := 0; i < 5; i++ {
		env.seedOrder(t, "cust_1", "pending")
	}

	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{
		"page": 1, "per_page": 2,
	})
	data := resp.data()
	orders, _ := data["orders"].([]any)
	if len(orders) != 2 {
		t.Fatalf("expected page size 2, got %d", len(orders))
	}
	if total, _ := data["total"].(float64); total != 5 {
		t.Fatalf("expected total 5, got %v", data["total"])
	}
}

func TestLoadOrdersInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "load-orders-ajax", "{oops")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
