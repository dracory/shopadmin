package discount_manager_test

import (
	"context"
	"testing"
)

func TestDiscountDeleteSuccess(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "C1", "active")

	resp := postAction(t, env, "delete-discount", nil, map[string]any{"discount_id": d.GetID()})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	found, _ := env.store.DiscountFindByID(context.Background(), d.GetID())
	if found != nil {
		t.Fatal("expected discount deleted")
	}
}

func TestDiscountDeleteMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-discount", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
