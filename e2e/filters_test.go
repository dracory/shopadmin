package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// openFiltersModal navigates to a manager page, waits for the Vue app to
// mount, and opens the Filter Conditions modal.
func openFiltersModal(t *testing.T, url string) context.Context {
	t.Helper()
	ctx := newTab(t)
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Poll(`!!document.querySelector('#app[data-v-app]')`, nil,
			chromedp.WithPollingInterval(200*time.Millisecond)),
		// Wait for the initial data load to finish so the filter bar is rendered.
		chromedp.Poll(
			`!!document.querySelector('#app table') || document.body.innerText.includes('No ')`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		// Click the "Filters" button in the filter bar.
		chromedp.Evaluate(`(function(){
			const btn = [...document.querySelectorAll('button')]
				.find(b => b.textContent.trim().startsWith('Filters'));
			if (!btn) return false;
			btn.click();
			return true;
		})()`, nil),
		chromedp.Poll(
			`document.body.innerText.includes('Filter Conditions')`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
	)
	if err != nil {
		t.Fatalf("page %s: filter modal did not open: %v", url, err)
	}
	return ctx
}

// applyStatusCondition selects the given field in the first condition
// row, picks a value, applies the filters, and returns the page text.
func applyCondition(t *testing.T, ctx context.Context, field, value string) string {
	t.Helper()
	var body string
	err := chromedp.Run(ctx,
		// Select the field in the first condition row.
		chromedp.Evaluate(`(function(){
			const row = document.querySelector('.modal.show .row');
			const sel = row.querySelector('select');
			sel.value = '`+field+`';
			sel.dispatchEvent(new Event('change', { bubbles: true }));
			return sel.value;
		})()`, nil),
		// Wait for the value input/select to render (3rd element in the row).
		chromedp.Poll(
			`document.querySelectorAll('.modal.show .row select, .modal.show .row input').length >= 3`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		// Set the value on the last form control of the first row.
		chromedp.Evaluate(`(function(){
			const row = document.querySelector('.modal.show .row');
			const ctrls = row.querySelectorAll('select, input');
			const val = ctrls[ctrls.length - 1];
			val.value = '`+value+`';
			val.dispatchEvent(new Event('change', { bubbles: true }));
			val.dispatchEvent(new Event('input', { bubbles: true }));
			return val.value;
		})()`, nil),
		// Click "Apply Filters".
		chromedp.Evaluate(`(function(){
			const btn = [...document.querySelectorAll('.modal.show button')]
				.find(b => b.textContent.trim() === 'Apply Filters');
			if (!btn) return false;
			btn.click();
			return true;
		})()`, nil),
		// Wait for the filter status text to reflect one active filter.
		chromedp.Poll(
			`document.body.innerText.includes('1 filter(s) active')`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("applying condition %s=%s failed: %v", field, value, err)
	}
	return body
}

func TestCategoryManagerFilterModal(t *testing.T) {
	ctx := openFiltersModal(t, adminURL("categories"))
	body := applyCondition(t, ctx, "status", "active")
	if !strings.Contains(body, "1 filter(s) active") {
		t.Fatal("expected '1 filter(s) active' after applying status filter")
	}
}

func TestProductManagerFilterModal(t *testing.T) {
	ctx := openFiltersModal(t, adminURL("products"))
	body := applyCondition(t, ctx, "status", "active")
	if !strings.Contains(body, "1 filter(s) active") {
		t.Fatal("expected '1 filter(s) active' after applying status filter")
	}
}

func TestDiscountManagerFilterModal(t *testing.T) {
	ctx := openFiltersModal(t, adminURL("discounts"))
	body := applyCondition(t, ctx, "status", "active")
	if !strings.Contains(body, "1 filter(s) active") {
		t.Fatal("expected '1 filter(s) active' after applying status filter")
	}
}

func TestOrderManagerFilterModal(t *testing.T) {
	ctx := openFiltersModal(t, adminURL("orders"))
	body := applyCondition(t, ctx, "status", "pending")
	if !strings.Contains(body, "1 filter(s) active") {
		t.Fatal("expected '1 filter(s) active' after applying status filter")
	}
}


