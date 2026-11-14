package order_details_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// escapeJSString is unexported; it is exercised through the page
// render path — an order ID containing a single quote must not break
// out of the embedded single-quoted JS string literal.
func TestOrderDetailsPageEscapesOrderID(t *testing.T) {
	env := newTestEnv(t)

	malicious := "x'onload='alert(1)"
	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=order_details&order_id="+malicious, http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.OrderDetails(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, "'x'onload='") {
		t.Fatal("expected order ID to be JS-escaped in page")
	}
	if !strings.Contains(body, `x\'onload=\'alert(1)`) && !strings.Contains(body, "x&#39;onload") {
		t.Fatalf("expected escaped order ID in page, got %q", body[:400])
	}
}
