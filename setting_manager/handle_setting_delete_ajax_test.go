package setting_manager_test

import "testing"

func TestSettingDeleteMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-setting", nil, map[string]any{})
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestSettingDeleteNotFound(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-setting", nil, map[string]any{"setting_id": "missing"})
	// Deleting a non-existent ID may succeed or error depending on
	// store semantics; just ensure a valid envelope is returned.
	if resp.Status == "" {
		t.Fatalf("expected status field, got %+v", resp)
	}
}
