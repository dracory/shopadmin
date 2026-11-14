package category_manager_test

import (
	"context"
	"testing"
)

func TestCategoryDeleteSuccess(t *testing.T) {
	env := newTestEnv(t)
	cat := env.seedCategory(t, "Doomed", "active")

	resp := postAction(t, env, "delete-category", nil, map[string]any{"category_id": cat.GetID()})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	found, err := env.store.CategoryFindByID(context.Background(), cat.GetID())
	if err != nil || found != nil {
		t.Fatal("expected category deleted")
	}
}

func TestCategoryDeleteMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-category", nil, map[string]any{})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
