package e2e

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestProductUpdateAssignsCategory verifies a category can be assigned
// to a product from the product update page: tick a checkbox in the
// Categories card, save, reload, and the tick must persist.
func TestProductUpdateAssignsCategory(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("product_update") + "&product_id=" + seedProductID + "&view=details"

	// The [id^="cat-"] selector only matches after Vue mounts and binds
	// :id — the raw template checkbox exists in the DOM before then and
	// clicking it is a no-op that Vue overwrites on mount.
	firstCheckbox := `#categories-app input.form-check-input[id^="cat-"]`
	saveButton := `#categories-app button.btn-success`

	// Load the page, wait for the categories checklist (AJAX), tick the
	// first category, and save.
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady(firstCheckbox, chromedp.ByQuery),
		chromedp.Click(firstCheckbox, chromedp.ByQuery),
		chromedp.Click(saveButton, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("categories card interaction failed: %v", err)
	}
	// Give the save fetch a moment before reloading.
	err = chromedp.Run(ctx, chromedp.Sleep(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	// Reload and assert a checkbox stayed checked.
	var checked int
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady(firstCheckbox, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelectorAll('#categories-app input.form-check-input:checked').length`, &checked),
	)
	if err != nil {
		t.Fatalf("reload of product update failed: %v", err)
	}
	if checked < 1 {
		t.Fatal("expected at least one category to remain checked after save + reload")
	}
}
