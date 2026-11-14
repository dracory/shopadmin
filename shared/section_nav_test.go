package shared

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/shopadmin/testutils"

	_ "modernc.org/sqlite"
)

func navRequest(controller, fileManagerURL string) *http.Request {
	r := httptest.NewRequest("GET", "/?controller="+controller, nil)
	ctx := context.WithValue(r.Context(), KeyShopAdminURL, "/admin/shop")
	if fileManagerURL != "" {
		ctx = context.WithValue(ctx, KeyFileManagerURL, fileManagerURL)
	}
	return r.WithContext(ctx)
}

func TestSectionNavRendersTabs(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	html := SectionNav(navRequest("", ""), store, nil).ToHTML()
	for _, want := range []string{"Dashboard", "Orders", "Products", "Categories", "Discounts", "Settings", "nav-pills"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in nav, got %q", want, html)
		}
	}
	if strings.Contains(html, "File Manager") {
		t.Fatal("expected no File Manager tab without configured URL")
	}
}

func TestSectionNavActiveTab(t *testing.T) {
	cases := map[string]string{
		"":                 "dashboard",
		CONTROLLER_ORDERS:  "orders",
		CONTROLLER_PRODUCTS: "products",
		CONTROLLER_CATEGORIES: "categories",
		CONTROLLER_DISCOUNTS:  "discounts",
		CONTROLLER_SETTINGS:   "settings",
	}
	for controller, _ := range cases {
		html := SectionNav(navRequest(controller, ""), nil, nil).ToHTML()
		if !strings.Contains(html, "active") {
			t.Fatalf("expected an active tab for controller %q", controller)
		}
	}
}

func TestSectionNavFileManagerTab(t *testing.T) {
	html := SectionNav(navRequest("", "/admin/files"), nil, nil).ToHTML()
	if !strings.Contains(html, "File Manager") || !strings.Contains(html, "/admin/files") {
		t.Fatalf("expected File Manager tab, got %q", html)
	}
}
