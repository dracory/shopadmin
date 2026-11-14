package product_view_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// escapeJSString is unexported; it is exercised through the page
// render path — a product ID containing a single quote must not
// break out of the embedded JS string literal.
func TestProductViewPageEscapesProductID(t *testing.T) {
	env := newTestEnv(t)

	malicious := "x'onload='alert(1)"
	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=product_view&product_id="+malicious, http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductView(rec, req)

	if strings.Contains(rec.Body.String(), "'x'onload='") {
		t.Fatal("expected product ID to be JS-escaped in page")
	}
}
