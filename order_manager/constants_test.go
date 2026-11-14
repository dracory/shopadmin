package order_manager_test

import "testing"

// The condition-filtering logic in constants.go and the
// CustomerResolver path are exercised through load-orders-ajax.

func TestLoadOrdersConditionsStatus(t *testing.T) {
	env := newTestEnv(t)
	env.seedOrder(t, "cust_1", "pending")
	env.seedOrder(t, "cust_2", "completed")

	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "status", "operator": "equals", "value": "completed"},
		},
	})
	orders, _ := resp.data()["orders"].([]any)
	if len(orders) != 1 {
		t.Fatalf("expected 1 completed order, got %d", len(orders))
	}
}

func TestLoadOrdersCustomerNameFilter(t *testing.T) {
	resolver := &mockCustomerResolver{ids: []string{"cust_9"}}
	env := newTestEnvWithResolver(t, resolver)
	env.seedOrder(t, "cust_9", "pending")
	env.seedOrder(t, "cust_other", "pending")

	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{
		"customer_name": "John",
	})
	orders, _ := resp.data()["orders"].([]any)
	if len(orders) != 1 {
		t.Fatalf("expected 1 order for resolved customer, got %d", len(orders))
	}
}

func TestLoadOrdersCustomerFilterWithoutResolver(t *testing.T) {
	env := newTestEnv(t)
	env.seedOrder(t, "cust_1", "pending")
	env.seedOrder(t, "cust_2", "pending")

	// Nil resolver: customer filters are skipped, all orders returned.
	resp := postAction(t, env, "load-orders-ajax", nil, map[string]any{
		"customer_name": "John",
	})
	orders, _ := resp.data()["orders"].([]any)
	if len(orders) != 2 {
		t.Fatalf("expected customer filter skipped, got %d orders", len(orders))
	}
}
