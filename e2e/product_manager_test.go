package e2e

import (
	"strings"
	"testing"
)

func TestProductManagerLoadsSeededProducts(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("products"))
	if !strings.Contains(body, "Product 01") {
		t.Fatalf("expected seeded 'Product 01' in product manager")
	}
}
