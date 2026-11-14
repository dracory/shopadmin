package category_manager_test

import (
	"context"
	"testing"
)

func TestCategoryDeleteSelectedSuccess(t *testing.T) {
	env := newTestEnv(t)
	a := env.seedCategory(t, "A", "active")
	b := env.seedCategory(t, "B", "active")
	c := env.seedCategory(t, "C", "active")

	resp := postAction(t, env, "delete-selected", nil, map[string]any{
		"bulk_category_ids": []string{a.GetID(), b.GetID()},
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	ctx := context.Background()
	if found, _ := env.store.CategoryFindByID(ctx, a.GetID()); found != nil {
		t.Fatal("expected A deleted")
	}
	if found, _ := env.store.CategoryFindByID(ctx, c.GetID()); found == nil {
		t.Fatal("expected C kept")
	}
}

func TestCategoryDeleteSelectedEmpty(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-selected", nil, map[string]any{"bulk_category_ids": []string{}})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
