package e2e

import (
	"strings"
	"testing"
)

func TestDiscountManagerLoadsSeededDiscounts(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("discounts"))
	if !strings.Contains(body, "DISC001") {
		t.Fatalf("expected seeded 'DISC001' in discount manager")
	}
}
