package category_manager_test

import "testing"

// applyConditions is unexported; its behavior is exercised through the
// load-categories action with conditions in the request body.

func TestLoadCategoriesConditions(t *testing.T) {
	env := newTestEnv(t)
	env.seedCategory(t, "Electronics", "active")
	env.seedCategory(t, "Books", "active")
	env.seedCategory(t, "Old Books", "draft")

	// search condition filters by title
	resp := postAction(t, env, "load-categories", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "search", "operator": "=", "value": "Book"},
		},
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	cats, _ := resp.data()["categories"].([]any)
	if len(cats) != 2 {
		t.Fatalf("expected 2 categories matching 'Book', got %d", len(cats))
	}

	// status condition
	resp = postAction(t, env, "load-categories", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "status", "operator": "=", "value": "draft"},
		},
	})
	cats, _ = resp.data()["categories"].([]any)
	if len(cats) != 1 {
		t.Fatalf("expected 1 draft category, got %d", len(cats))
	}

	// empty value conditions are ignored
	resp = postAction(t, env, "load-categories", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "status", "operator": "=", "value": "  "},
		},
	})
	cats, _ = resp.data()["categories"].([]any)
	if len(cats) != 3 {
		t.Fatalf("expected all 3 categories, got %d", len(cats))
	}
}
