package e2e

import (
	"context"
	"strings"
	"testing"

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
