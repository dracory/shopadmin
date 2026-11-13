package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestSettingManagerLoads verifies the settings page mounts its Vue app
// and finishes loading (table or empty-state alert).
func TestSettingManagerLoads(t *testing.T) {
	ctx := newTab(t)

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("settings")),
		// data-v-app is stamped by Vue on mount — DOM selectors alone
		// can match the uncompiled template.
		chromedp.Poll(`!!document.querySelector('#app[data-v-app]')`, nil,
			chromedp.WithPollingInterval(200*time.Millisecond)),
		// Initial settings fetch finished when either the table or the
		// "No settings found" empty state renders.
		chromedp.Poll(
			`!!document.querySelector('#app table') || document.body.innerText.includes('No settings found')`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("settings page failed: %v", err)
	}
	if !strings.Contains(body, "New Setting") {
		t.Fatal("expected 'New Setting' button on settings page")
	}
}
