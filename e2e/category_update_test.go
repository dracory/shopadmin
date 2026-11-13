package e2e

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestCategoryUpdateFlow verifies the category update page loads the
// existing category into the form and that saving persists the change.
func TestCategoryUpdateFlow(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("category_update") + "&category_id=" + seedCategoryID
	newTitle := "E2E Updated " + time.Now().Format("150405")

	// The form prefills via the load-category fetch — wait for a
	// non-empty title input value, not just the element.
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#category-title", chromedp.ByID),
		chromedp.Poll(`document.querySelector('#category-title').value !== ''`, nil,
			chromedp.WithPollingInterval(200*time.Millisecond)),
	)
	if err != nil {
		t.Fatalf("category update form did not prefill: %v", err)
	}

	// Change the title and save — the JS redirects back to the manager.
	err = chromedp.Run(ctx,
		chromedp.SetValue("#category-title", newTitle, chromedp.ByID),
		chromedp.Click("#btn-update-category", chromedp.ByID),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		t.Fatalf("category update submit failed: %v", err)
	}

	// Reopen the page; the new title must persist.
	var persisted string
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#category-title", chromedp.ByID),
		chromedp.Poll(`document.querySelector('#category-title').value !== ''`, nil,
			chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Value("#category-title", &persisted, chromedp.ByID),
	)
	if err != nil {
		t.Fatalf("reload of category update failed: %v", err)
	}
	if persisted != newTitle {
		t.Fatalf("expected persisted title %q, got %q", newTitle, persisted)
	}
}
