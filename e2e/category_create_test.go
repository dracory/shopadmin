package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestCategoryCreateFlow verifies the "New Category" button exists on
// the categories page and that clicking through lets a user create a
// category that then appears in the manager.
func TestCategoryCreateFlow(t *testing.T) {
	ctx := newTab(t)
	title := "E2E Category " + time.Now().Format("150405")

	// The button must exist on the categories page. Waiting for the
	// second rendered table row also guarantees Vue has mounted and the
	// button's :href is bound — clicking earlier hits an empty href.
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("categories")),
		chromedp.WaitReady("table tbody tr:nth-of-type(2)", chromedp.ByQuery),
		chromedp.WaitReady("#btn-new-category", chromedp.ByID),
	)
	if err != nil {
		t.Fatalf("New Category button missing on categories page: %v", err)
	}

	// Click it → land on the create form.
	err = chromedp.Run(ctx,
		chromedp.Click("#btn-new-category", chromedp.ByID),
		chromedp.WaitReady("#category-title", chromedp.ByID),
	)
	if err != nil {
		t.Fatalf("create form did not load: %v", err)
	}

	// Fill and submit.
	err = chromedp.Run(ctx,
		chromedp.SendKeys("#category-title", title, chromedp.ByID),
		chromedp.Click("#btn-create-category", chromedp.ByID),
	)
	if err != nil {
		t.Fatalf("create form interaction failed: %v", err)
	}

	// The JS redirects back to the category manager; the new category
	// must appear in the table.
	var body string
	err = chromedp.Run(ctx,
		chromedp.WaitReady("table tbody tr:nth-of-type(2)", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("redirect back to categories failed: %v", err)
	}
	if !strings.Contains(body, title) {
		t.Fatalf("expected newly created category %q in manager", title)
	}
}
