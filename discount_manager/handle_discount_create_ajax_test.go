package discount_manager_test

import (
	"context"
	"strings"
	"testing"
)

func TestDiscountCreateSuccess(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "create-discount-ajax", nil, map[string]any{"title": "Spring Sale"})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	id, _ := resp.data()["discount_id"].(string)
	if id == "" {
		t.Fatal("expected discount_id in response")
	}
	discount, err := env.store.DiscountFindByID(context.Background(), id)
	if err != nil || discount == nil || discount.GetTitle() != "Spring Sale" {
		t.Fatalf("expected persisted discount, got %v", err)
	}
}

func TestDiscountCreateMissingTitle(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "create-discount-ajax", nil, map[string]any{"title": "  "})
	if resp.status() != "error" || !strings.Contains(resp.message(), "Title") {
		t.Fatalf("expected title error, got %+v", resp)
	}
}

func TestDiscountCreateInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "create-discount-ajax", "{bad")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
