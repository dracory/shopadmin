package discount_update_test

import (
	"context"
	"testing"
)

// The details component is exercised through the controller's
// load-details and save-details actions.

func TestLoadDetails(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "SAVE10", "active")
	d.SetDescription("ten off")
	_ = env.store.DiscountUpdate(context.Background(), d)

	resp := postAction(t, env, "load-details", d.GetID(), nil)
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if data["code"] != "SAVE10" || data["status"] != "active" {
		t.Fatalf("unexpected loaded data %+v", data)
	}
}

func TestLoadDetailsNotFound(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "load-details", "missing", nil)
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestSaveDetails(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "SAVE10", "active")

	resp := postAction(t, env, "save-details", d.GetID(), map[string]any{
		"code":        "SAVE20",
		"title":       "Twenty off",
		"type":        "percent",
		"amount":      20,
		"status":      "active",
		"description": "updated",
	})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	updated, err := env.store.DiscountFindByID(context.Background(), d.GetID())
	if err != nil || updated == nil {
		t.Fatalf("find updated: %v", err)
	}
	if updated.GetCode() != "SAVE20" || updated.GetTitle() != "Twenty off" {
		t.Fatalf("unexpected updated fields: %q %q", updated.GetCode(), updated.GetTitle())
	}
}

func TestSaveDetailsDuplicateCode(t *testing.T) {
	env := newTestEnv(t)
	_ = env.seedDiscount(t, "TAKEN", "active")
	d := env.seedDiscount(t, "SAVE10", "active")

	resp := postAction(t, env, "save-details", d.GetID(), map[string]any{
		"code":   "TAKEN",
		"title":  "x",
		"type":   "percent",
		"amount": 5,
		"status": "active",
	})
	if resp.status() != "error" {
		t.Fatalf("expected duplicate-code error, got %+v", resp)
	}
}

func TestSaveDetailsNotFound(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "save-details", "missing", map[string]any{"code": "X"})
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
