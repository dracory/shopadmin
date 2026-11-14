package setting_manager_test

import (
	"context"
	"testing"
)

func TestLoadSettingsPagination(t *testing.T) {
	env := newTestEnv(t)
	for _, k := range []string{"k1", "k2", "k3"} {
		if err := env.store.SettingSet(context.Background(), k, "v"); err != nil {
			t.Fatalf("seed setting: %v", err)
		}
	}

	resp := postAction(t, env, "load-settings", nil, map[string]any{
		"page": 0, "per_page": 2,
	})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	settings, _ := resp.Data["settings"].([]any)
	if len(settings) != 2 {
		t.Fatalf("expected page size 2, got %d", len(settings))
	}
	if total, _ := resp.Data["total"].(float64); total < 3 {
		t.Fatalf("expected total >= 3, got %v", resp.Data["total"])
	}
}

func TestLoadSettingsInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRawBody(t, env, "load-settings", "{oops")
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
