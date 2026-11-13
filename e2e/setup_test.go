// Package e2e contains end-to-end tests that drive the shopadmin panel
// with a real headless Chrome via chromedp (CDP). The admin is served
// in-process by httptest on a random port — no binary build needed.
package e2e

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dracory/shopadmin"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"
	_ "modernc.org/sqlite"
)

var (
	baseURL       string
	allocCtx      context.Context
	testStore     shopstore.StoreInterface
	seedProductID string
)

const adminPath = "/admin/shop"

func TestMain(m *testing.M) {
	// In-process admin server backed by a seeded in-memory SQLite DB.
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		panic("e2e: failed to init store: " + err.Error())
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	testutils.SeedDB(store, logger)
	testStore = store

	products, err := store.ProductList(context.Background(), shopstore.NewProductQuery())
	if err != nil || len(products) == 0 {
		panic("e2e: no seeded products: " + fmt.Sprint(err))
	}
	seedProductID = products[0].GetID()

	admin, err := shopadmin.New(shopadmin.AdminOptions{
		Store:        store,
		Logger:       logger,
		AdminHomeURL: "/admin",
		ShopAdminURL: adminPath,
		// AuthUserID intentionally nil — the panel is open.
	})
	if err != nil {
		panic("e2e: failed to create shopadmin: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle(adminPath, http.HandlerFunc(admin.Handle))
	mux.Handle(adminPath+"/", http.HandlerFunc(admin.Handle))
	srv := httptest.NewServer(mux)
	baseURL = srv.URL

	// One shared browser process for the entire test run.
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Headless)
	if bin := chromeBin(); bin != "" {
		opts = append(opts, chromedp.ExecPath(bin))
	}
	var cancel context.CancelFunc
	allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)

	code := m.Run()

	cancel()
	srv.Close()
	os.Exit(code)
}

// newTab opens a fresh tab in the shared browser and returns its
// context. Each tab gets a 30s deadline; cleanup closes the tab.
func newTab(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := chromedp.NewContext(allocCtx)
	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	t.Cleanup(func() { cancelTimeout(); cancel() })
	return ctx
}

// chromeBin resolves a headless Chromium binary. Order: CHROME_BIN env,
// Playwright's downloaded headless shell, empty (chromedp finds Chrome
// on PATH / standard install locations).
func chromeBin() string {
	if b := os.Getenv("CHROME_BIN"); b != "" {
		return b
	}
	patterns := []string{
		// Windows
		filepath.Join(os.Getenv("LOCALAPPDATA"),
			"ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe"),
		// WSL: Windows-side Playwright shell via interop
		"/mnt/*/Users/*/AppData/Local/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe",
		// Linux
		filepath.Join(os.Getenv("HOME"),
			".cache/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-linux64/chrome-headless-shell"),
	}
	for _, pattern := range patterns {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	return ""
}
