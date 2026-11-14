package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dracory/shopstore"
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
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady(firstCheckbox, chromedp.ByQuery),
		chromedp.Click(firstCheckbox, chromedp.ByQuery),
		chromedp.Click(saveButton, chromedp.ByQuery),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("categories card interaction failed: %v", err)
	}
	if !strings.Contains(body, "Shop. Edit Product") {
		t.Fatal("expected 'Shop. Edit Product' header on product update page")
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

// TestProductUpdateAddsVariantDimension verifies the Variants card can add
// a dimension, edit the generated price matrix, save it, and reload it.
func TestProductUpdateAddsVariantDimension(t *testing.T) {
	product := shopstore.NewProduct()
	product.SetTitle("Variant Matrix Product")
	product.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)
	if err := testStore.ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	ctx := newTab(t)
	url := adminURL("product_update") + "&product_id=" + product.GetID() + "&view=details"

	var rowCount int
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady(`#variants-app .alert`, chromedp.ByQuery),
		chromedp.Click(`#variants-app fieldset button.btn-primary`, chromedp.ByQuery),
		chromedp.WaitVisible(`#variants-app .modal.show`, chromedp.ByQuery),
		chromedp.SetValue(`#variants-app .modal.show .mb-3 input`, "size", chromedp.ByQuery),
		chromedp.SetValue(`#variants-app .modal.show .modal-body > div:last-child input`, "S, M", chromedp.ByQuery),
		chromedp.Click(`#variants-app .modal.show .modal-footer .btn-primary`, chromedp.ByQuery),
		chromedp.WaitVisible(`#variants-app table tbody input[placeholder="SKU"]`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelectorAll('#variants-app table tbody tr').length`, &rowCount),
	)
	if err != nil {
		t.Fatalf("add dimension interaction failed: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("expected 2 generated variant rows, got %d", rowCount)
	}

	err = chromedp.Run(ctx,
		chromedp.SendKeys(`#variants-app table tbody tr:nth-of-type(1) input[placeholder="SKU"]`, "MATRIX-S", chromedp.ByQuery),
		chromedp.SendKeys(`#variants-app table tbody tr:nth-of-type(1) input[placeholder="parent price"]`, "12.34", chromedp.ByQuery),
		chromedp.SendKeys(`#variants-app table tbody tr:nth-of-type(1) td:nth-of-type(4) input`, "7", chromedp.ByQuery),
		chromedp.Click(`#variants-app .card-header .btn-success`, chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#variants-app table tbody input[placeholder="SKU"]`, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("price matrix save/reload failed: %v", err)
	}

	var values []string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('#variants-app table tbody tr:nth-of-type(1) input')).map(input => input.value)`, &values),
	)
	if err != nil {
		t.Fatalf("failed to read reloaded matrix values: %v", err)
	}
	if len(values) != 3 || values[0] != "MATRIX-S" || values[1] != "12.34" || values[2] != "7" {
		t.Fatalf("expected persisted matrix values [MATRIX-S 12.34 7], got %v", values)
	}
}
