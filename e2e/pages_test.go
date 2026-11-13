package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// adminURL builds a controller URL, e.g. /admin/shop?controller=products.
func adminURL(controller string) string {
	return baseURL + adminPath + "?controller=" + controller
}

// waitTable navigates to a manager page and waits until Vue has loaded
// data and rendered the results table. Returns the page's text content.
func waitTable(t *testing.T, ctx context.Context, url string) string {
	t.Helper()
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("table", chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("page %s failed: %v", url, err)
	}
	return body
}

func TestProductManagerLoadsSeededProducts(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("products"))
	if !strings.Contains(body, "Product 01") {
		t.Fatalf("expected seeded 'Product 01' in product manager")
	}
}

func TestCategoryManagerLoadsSeededCategories(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("categories"))
	if !strings.Contains(body, "Category 01") {
		t.Fatalf("expected seeded 'Category 01' in category manager")
	}
}

func TestDiscountManagerLoadsSeededDiscounts(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("discounts"))
	if !strings.Contains(body, "DISC001") {
		t.Fatalf("expected seeded 'DISC001' in discount manager")
	}
}

func TestOrderManagerLoadsSeededOrders(t *testing.T) {
	body := waitTable(t, newTab(t), adminURL("orders"))
	// Customer names resolve via CustomerResolver, which is nil in e2e;
	// assert on the seeded order statuses instead.
	if !strings.Contains(body, "pending") {
		t.Fatalf("expected seeded order with status 'pending' in order manager")
	}
}

// TestCategoryCreateFlow verifies the "New Category" button exists on
// the categories page and that clicking through lets a user create a
// category that then appears in the manager.
func TestCategoryCreateFlow(t *testing.T) {
	ctx := newTab(t)
	title := "E2E Category " + time.Now().Format("150405")

	// The button must exist on the categories page. Waiting for the
	// Vue-rendered table also guarantees the app has mounted and the
	// button's :href is bound — clicking earlier hits an empty href.
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("categories")),
		chromedp.WaitReady("table", chromedp.ByQuery),
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
		chromedp.WaitReady("table", chromedp.ByQuery),
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

// TestProductUpdateAssignsCategory verifies a category can be assigned
// to a product from the product update page: tick a checkbox in the
// Categories card, save, reload, and the tick must persist.
func TestProductUpdateAssignsCategory(t *testing.T) {
	ctx := newTab(t)
	url := adminURL("product_update") + "&product_id=" + seedProductID + "&view=details"

	firstCheckbox := `#categories-app input.form-check-input`
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

func TestHomePageLoads(t *testing.T) {
	ctx := newTab(t)
	var body string
	err := chromedp.Run(ctx, chromedp.Navigate(baseURL+adminPath))
	if err != nil {
		t.Fatalf("navigate: %v", err)
	}
	err = chromedp.Run(ctx,
		chromedp.WaitReady("body"),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if strings.TrimSpace(body) == "" {
		t.Fatal("expected non-empty home page")
	}
}
