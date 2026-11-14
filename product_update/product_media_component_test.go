package product_update_test

import (
	"net/url"
	"testing"
)

// The media component is exercised through the controller's
// load-media action (upload requires multipart form data).

func TestLoadMediaEmpty(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	resp := postAction(t, env, "load-media",
		url.Values{"product_id": {p.GetID()}}, nil)
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}
