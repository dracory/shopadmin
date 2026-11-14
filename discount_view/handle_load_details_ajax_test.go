package discount_view_test

import "testing"

func TestLoadDetailsSuccess(t *testing.T) {
	env := newTestEnv(t)
	d := env.seedDiscount(t, "VIEW10", "active")

	resp := getAction(t, env, "action=load-details&discount_id="+d.GetID())
	if resp.status() != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
	data := resp.data()
	if data["code"] != "VIEW10" {
		t.Fatalf("unexpected data %+v", data)
	}
	if _, ok := data["amount_formatted"]; !ok {
		t.Fatal("expected amount_formatted in response")
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

	resp := getAction(t, env, "action=load-details&discount_id=missing")
	if resp.status() != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
