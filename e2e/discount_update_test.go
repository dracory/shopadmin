package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestDiscountUpdatePrefillsForm verifies the discount update page
// loads the seeded discount into the Vue form (v-model fields populate
// only after mount + AJAX load completes).
func TestDiscountUpdatePrefillsForm(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("discount_update") + "&discount_id=" + seedDiscountID

	var title string
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// data-v-app is stamped by Vue on mount; the disabled ID input
		// only gets a value once the details AJAX load completes.
		chromedp.Poll(`!!document.querySelector('#details-app[data-v-app]')`, nil,
			chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Poll(
			`(document.querySelector('#details-app input[disabled]')||{}).value !== ''`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Value(`#details-app input[placeholder="Enter discount title..."]`, &title, chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("discount update page failed: %v", err)
	}
	if !strings.Contains(body, "Shop. Edit Discount") {
		t.Fatal("expected 'Shop. Edit Discount' header on discount update page")
	}
	if title != seedDiscountTitle {
		t.Fatalf("expected prefilled title %q, got %q", seedDiscountTitle, title)
	}
}
