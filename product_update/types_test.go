package product_update_test

import (
	"net/http/httptest"
	"testing"

	"github.com/dracory/shopadmin/product_update"
)

func TestReqArrayOfMaps(t *testing.T) {
	req := httptest.NewRequest("GET",
		"/?tags[0][name]=color&tags[0][value]=red&tags[1][name]=size&other=1", nil)

	got := product_update.ReqArrayOfMaps(req, "tags", nil)
	if len(got) != 2 {
		t.Fatalf("expected 2 maps, got %d (%v)", len(got), got)
	}
	seen := map[string]bool{}
	for _, m := range got {
		seen[m["name"]] = true
	}
	if !seen["color"] || !seen["size"] {
		t.Fatalf("unexpected maps %v", got)
	}
}

func TestReqArrayOfMapsEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/?x=1", nil)
	if got := product_update.ReqArrayOfMaps(req, "tags", nil); len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
