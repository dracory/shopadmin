package setting_manager_test

import "testing"

func TestSettingSaveMissingKey(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-setting", nil, map[string]any{"key": "", "value": "v"})
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestSettingSaveInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRawBody(t, env, "save-setting", "{oops")
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
