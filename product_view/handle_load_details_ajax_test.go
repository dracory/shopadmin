package product_view_test

import "testing"

func TestLoadDetailsSuccess(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t, "Viewable", "active")

	resp := getAction(t, env, "action=load-details&product_id="+p.GetID())
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if data["title"] != "Viewable" {
		t.Fatalf("unexpected data %+v", data)
	}
	if _, ok := data["price_formatted"]; !ok {
		t.Fatal("expected price_formatted in response")
	}
}

func TestLoadDetailsMissingID(t *testing.T) {
	env := newTestEnv(t)

	resp := getAction(t, env, "action=load-details")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestLoadDetailsNotFound(t *testing.T) {
	env := newTestEnv(t)

	resp := getAction(t, env, "action=load-details&product_id=missing")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
