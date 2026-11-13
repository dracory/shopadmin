package product_update_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/dracory/shopstore"
)

// TestSaveDetailsAcceptsStringPrice verifies the standard payload where
// price and quantity arrive as JSON strings.
func TestSaveDetailsAcceptsStringPrice(t *testing.T) {
	env := newTestEnv(t)
	product := env.seedProduct(t)

	q := url.Values{"product_id": {product.GetID()}}
	resp := postAction(t, env, "save-details", q, map[string]any{
		"status":      shopstore.PRODUCT_STATUS_ACTIVE,
		"title":       "Updated Title",
		"description": "new description",
		"price":       "19.99",
		"quantity":    "7",
		"memo":        "new memo",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.ProductFindByID(context.Background(), product.GetID())
	if saved == nil {
		t.Fatal("expected saved product, got nil")
	}
	if saved.GetPrice() != "19.99" || saved.GetQuantity() != "7" {
		t.Errorf("expected price=19.99 quantity=7, got price=%q quantity=%q", saved.GetPrice(), saved.GetQuantity())
	}
	if saved.GetTitle() != "Updated Title" {
		t.Errorf("expected title updated, got %q", saved.GetTitle())
	}
}

// TestSaveDetailsAcceptsNumericPrice is a regression test for the
// "Invalid request body" error seen when editing the price in the UI.
// Vue's v-model on <input type="number"> auto-casts the value to a
// number, so the browser sends {"price": 25} rather than {"price": "25"}.
// The handler must accept both representations.
func TestSaveDetailsAcceptsNumericPrice(t *testing.T) {
	env := newTestEnv(t)
	product := env.seedProduct(t)

	q := url.Values{"product_id": {product.GetID()}}
	resp := postAction(t, env, "save-details", q, map[string]any{
		"status":      shopstore.PRODUCT_STATUS_ACTIVE,
		"title":       "Updated Title",
		"description": "new description",
		"price":       24.5,
		"quantity":    3,
		"memo":        "new memo",
	})
	if resp.Status != "success" {
		t.Fatalf("expected status=success, got %q (message=%q)", resp.Status, resp.Message)
	}

	saved, _ := env.store.ProductFindByID(context.Background(), product.GetID())
	if saved == nil {
		t.Fatal("expected saved product, got nil")
	}
	if saved.GetPrice() != "24.5" || saved.GetQuantity() != "3" {
		t.Errorf("expected price=24.5 quantity=3, got price=%q quantity=%q", saved.GetPrice(), saved.GetQuantity())
	}
}
