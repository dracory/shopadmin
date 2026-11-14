package category_manager_test

import "testing"

func TestLoadCategoriesSuccess(t *testing.T) {
	env := newTestEnv(t)
	env.seedCategory(t, "Alpha", "active")
	env.seedCategory(t, "Beta", "draft")

	resp := postAction(t, env, "load-categories", nil, map[string]any{})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if total, _ := data["total"].(float64); total != 2 {
		t.Fatalf("expected total 2, got %v", data["total"])
	}
	cats, _ := data["categories"].([]any)
	if len(cats) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(cats))
	}
	if _, ok := data["urls"].(map[string]any); !ok {
		t.Fatal("expected urls in response")
	}
}

func TestLoadCategoriesStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	env.seedCategory(t, "Alpha", "active")
	env.seedCategory(t, "Beta", "draft")

	resp := postAction(t, env, "load-categories", nil, map[string]any{"status": "draft"})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	cats, _ := resp.data()["categories"].([]any)
	if len(cats) != 1 {
		t.Fatalf("expected 1 filtered category, got %d", len(cats))
	}
}

func TestLoadCategoriesInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "load-categories", "{bad json")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
