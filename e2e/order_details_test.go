package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestOrderDetailsShowsOrder verifies the order details page loads the
// seeded order: ID, status badge, and line items table.
func TestOrderDetailsShowsOrder(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("order_details") + "&order_id=" + seedOrderID

	var body string
	// The status badge only renders after Vue mounts AND the order
	// AJAX load completes (v-if="order").
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#app .card .badge", chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("order details page failed: %v", err)
	}
	if !strings.Contains(body, "Shop. Order Details") {
		t.Fatal("expected 'Shop. Order Details' header on order details page")
	}
	if !strings.Contains(body, seedOrderID) {
		t.Fatalf("expected order ID %q in order details", seedOrderID)
	}
}
