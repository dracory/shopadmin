package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestSettingManagerLoads verifies the settings page mounts its Vue app
// and finishes loading (settings editor or empty-state text).
func TestSettingManagerLoads(t *testing.T) {
	ctx := newTab(t)

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(adminURL("settings")),
		// data-v-app is stamped by Vue on mount — DOM selectors alone
		// can match the uncompiled template.
		chromedp.Poll(`!!document.querySelector('#app[data-v-app]')`, nil,
			chromedp.WithPollingInterval(200*time.Millisecond)),
		// Initial settings fetch finished when either the editor list or
		// the "No settings" empty state renders.
		chromedp.Poll(
			`!!document.querySelector('#app .list-group') || document.body.innerText.includes('No settings to edit')`,
			nil, chromedp.WithPollingInterval(200*time.Millisecond)),
		chromedp.Text("body", &body),
	)
	if err != nil {
		t.Fatalf("settings page failed: %v", err)
	}
	if !strings.Contains(body, "Shop. Settings") {
		t.Fatal("expected 'Shop. Settings' header on settings page")
	}
}
