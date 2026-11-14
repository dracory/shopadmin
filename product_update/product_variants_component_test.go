package product_update_test

import (
	"net/url"
	"testing"
)

// The variants component is exercised through the controller's
// load-variants and save-variants actions.

func TestSaveAndLoadVariants(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "save-variants",
		url.Values{"product_id": {p.GetID()}},
		map[string]any{
			"dimensions": []map[string]any{
				{"name": "Size", "required": true, "options": []string{"S", "M"}},
			},
			"variants": []map[string]any{
				{"values": map[string]string{"Size": "S"}, "sku": "W-S", "price": "9.99", "quantity": "3"},
			},
		})
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}

	resp = postAction(t, env, "load-variants",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
