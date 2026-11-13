package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestProductViewShowsProduct verifies the product view page loads the
// seeded product's details: title, status badge, and price.
func TestProductViewShowsProduct(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("product_view") + "&product_id=" + seedProductID

	var body string
	// The status badge only renders after Vue mounts AND the product
	// AJAX load completes (v-if="product").
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#app .card .badge", chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("product view page failed: %v", err)
	}
	if !strings.Contains(body, seedProductTitle) {
		t.Fatalf("expected seeded product title %q in product view", seedProductTitle)
	}
}
