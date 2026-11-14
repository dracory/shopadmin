package setting_manager_test

import (
	"context"
	"testing"
)

// applyConditions is unexported; it is exercised through the
// load-settings action with conditions in the request body.

func TestLoadSettingsConditions(t *testing.T) {
	env := newTestEnv(t)
	for _, k := range []string{"site_title", "site_logo", "mail_host"} {
		if err := env.store.SettingSet(context.Background(), k, "v"); err != nil {
			t.Fatalf("seed setting: %v", err)
		}
	}

	resp := postAction(t, env, "load-settings", nil, map[string]any{
		"conditions": []map[string]any{
			{"field": "key", "operator": "=", "value": "site"},
		},
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	settings, _ := resp.Data["settings"].([]any)
	if len(settings) != 2 {
		t.Fatalf("expected 2 settings matching 'site', got %d", len(settings))
	}
}
