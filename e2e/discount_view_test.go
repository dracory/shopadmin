package e2e

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestDiscountViewShowsDiscount verifies the discount view page loads
// the seeded discount's details: title and status badge.
func TestDiscountViewShowsDiscount(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("discount_view") + "&discount_id=" + seedDiscountID

	var body string
	// The status badge only renders after Vue mounts AND the discount
	// AJAX load completes (v-if="discount").
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#app .card .badge", chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("discount view page failed: %v", err)
	}
	if !strings.Contains(body, seedDiscountTitle) {
		t.Fatalf("expected seeded discount title %q in discount view", seedDiscountTitle)
	}
}
