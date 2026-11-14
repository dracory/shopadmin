package setting_manager_test

import "testing"

func TestSettingDeleteSelectedEmpty(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "delete-selected", nil, map[string]any{"bulk_setting_ids": []string{}})
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
