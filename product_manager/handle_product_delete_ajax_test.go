package product_manager_test

import (
	"context"
	"testing"
)

func TestProductDeleteSuccess(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t, "Doomed", "active")

	resp := postAction(t, env, "delete-product", nil, map[string]any{"product_id": p.GetID()})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	found, _ := env.store.ProductFindByID(context.Background(), p.GetID())
	if found != nil {
		t.Fatal("expected product deleted")
	}
}

func TestProductDeleteMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-product", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
