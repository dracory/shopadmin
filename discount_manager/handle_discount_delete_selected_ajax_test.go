package discount_manager_test

import (
	"context"
	"testing"
)

func TestDiscountDeleteSelectedSuccess(t *testing.T) {
	env := newTestEnv(t)
	a := env.seedDiscount(t, "A", "active")
	b := env.seedDiscount(t, "B", "active")
	c := env.seedDiscount(t, "C", "active")

	resp := postAction(t, env, "delete-selected", nil, map[string]any{
		"bulk_discount_ids": []string{a.GetID(), b.GetID()},
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	ctx := context.Background()
	if found, _ := env.store.DiscountFindByID(ctx, a.GetID()); found != nil {
		t.Fatal("expected A deleted")
	}
	if found, _ := env.store.DiscountFindByID(ctx, c.GetID()); found == nil {
		t.Fatal("expected C kept")
	}
}

func TestDiscountDeleteSelectedEmpty(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-selected", nil, map[string]any{"bulk_discount_ids": []string{}})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
