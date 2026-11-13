package e2e

import (
	"strings"
	"testing"
)

func TestOrderManagerLoadsSeededOrders(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("orders"))
	// Customer names resolve via CustomerResolver, which is nil in e2e;
	// assert on the seeded order statuses instead.
	if !strings.Contains(body, "pending") {
		t.Fatalf("expected seeded order with status 'pending' in order manager")
	}
}
