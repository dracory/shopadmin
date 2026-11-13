package setting_manager_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/dracory/shopstore"
)

func seedSetting(t *testing.T, env *testEnv, key, value string) shopstore.SettingInterface {
	t.Helper()
	s := shopstore.NewSetting().SetKey(key).SetValue(value)
	if err := env.store.SettingCreate(context.Background(), s); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return s
}

func TestSettingSaveAndLoad(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-setting", nil, map[string]any{
		"key":   "currency",
		"value": "GBP",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	resp = postAction(t, env, "load-settings", nil, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	settings, ok := resp.Data["settings"].([]interface{})
	if !ok || len(settings) != 1 {
		t.Fatalf("expected 1 setting in response, got %v", resp.Data["settings"])
	}

	if resp.Data["currency"] != "GBP" {
		t.Fatalf("expected currency=GBP, got %v", resp.Data["currency"])
	}
}

func TestSettingSaveUpdatesExistingKey(t *testing.T) {
	env := newTestEnv(t)

	postAction(t, env, "save-setting", nil, map[string]any{"key": "currency", "value": "GBP"})
	resp := postAction(t, env, "save-setting", nil, map[string]any{"key": "currency", "value": "EUR"})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q", resp.Status)
	}

	value, err := env.store.SettingGet(context.Background(), "currency", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "EUR" {
		t.Fatalf("expected EUR, got %q", value)
	}

	count, err := env.store.SettingCount(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 setting (upsert), got %d", count)
	}
}

func TestSettingSaveRequiresKey(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-setting", nil, map[string]any{
		"key":   "",
		"value": "x",
	})
	if resp.Status != "error" {
		t.Fatalf("expected status=error, got %q", resp.Status)
	}
}

func TestSettingDeleteAndDeleteSelected(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	s1 := seedSetting(t, env, "a", "1")
	s2 := seedSetting(t, env, "b", "2")

	resp := postAction(t, env, "delete-setting", nil, map[string]any{
		"setting_id": s1.GetID(),
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q", resp.Status)
	}

	resp = postAction(t, env, "delete-selected", nil, map[string]any{
		"bulk_setting_ids": []string{s2.GetID()},
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q", resp.Status)
	}

	count, _ := env.store.SettingCount(ctx, nil)
	if count != 0 {
		t.Fatalf("expected 0 settings, got %d", count)
	}
}

func TestLoadSettingsFiltersByKey(t *testing.T) {
	env := newTestEnv(t)

	for _, kv := range [][2]string{{"currency", "GBP"}, {"shop_name", "Test"}} {
		seedSetting(t, env, kv[0], kv[1])
	}

	resp := postAction(t, env, "load-settings", url.Values{}, map[string]any{
		"key": "curr",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q", resp.Status)
	}

	settings, ok := resp.Data["settings"].([]interface{})
	if !ok || len(settings) != 1 {
		t.Fatalf("expected 1 filtered setting, got %v", resp.Data["settings"])
	}
}
