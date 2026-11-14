package product_update_test

import (
	"net/url"
	"testing"
)

// The details component is exercised through the controller's
// load-details and save-details actions (save-details is covered in
// handle_save_details_ajax_test.go).

func TestLoadDetailsData(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "load-details",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.Data
	if data["price"] == nil || data["quantity"] == nil {
		t.Fatalf("expected price/quantity in details, got %+v", data)
	}
}
