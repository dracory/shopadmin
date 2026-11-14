package product_manager_test

import (
	"context"
	"testing"
)

func TestProductDeleteSelectedSuccess(t *testing.T) {
	env := newTestEnv(t)
	a := env.seedProduct(t, "A", "active")
	b := env.seedProduct(t, "B", "active")
	c := env.seedProduct(t, "C", "active")

	resp := postAction(t, env, "delete-selected", nil, map[string]any{
		"bulk_product_ids": []string{a.GetID(), b.GetID()},
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	ctx := context.Background()
	if found, _ := env.store.ProductFindByID(ctx, a.GetID()); found != nil {
		t.Fatal("expected A deleted")
	}
	if found, _ := env.store.ProductFindByID(ctx, c.GetID()); found == nil {
		t.Fatal("expected C kept")
	}
}

func TestProductDeleteSelectedEmpty(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-selected", nil, map[string]any{"bulk_product_ids": []string{}})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
