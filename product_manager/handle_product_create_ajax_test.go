package product_manager_test

import (
	"context"
	"strings"
	"testing"
)

func TestProductCreateSuccess(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "create-product-ajax", nil, map[string]any{"title": "New Widget"})
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	id, _ := resp.data()["product_id"].(string)
	if id == "" {
		t.Fatal("expected product_id in response")
	}
	product, err := env.store.ProductFindByID(context.Background(), id)
	if err != nil || product == nil || product.GetTitle() != "New Widget" {
		t.Fatalf("expected persisted product, got %v", err)
	}
}

func TestProductCreateMissingTitle(t *testing.T) {
	env := newTestEnv(t)

	resp := postAction(t, env, "create-product-ajax", nil, map[string]any{"title": "  "})
	if resp.status() != "error" || !strings.Contains(resp.message(), "Title") {
		t.Fatalf("expected title error, got %+v", resp)
	}
}

func TestProductCreateInvalidBody(t *testing.T) {
	env := newTestEnv(t)

	resp := postRaw(t, env, "create-product-ajax", "{bad")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
