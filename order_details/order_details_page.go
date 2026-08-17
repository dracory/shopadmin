package order_details

import (
	_ "embed"
	"html"
	"net/http"
	"strings"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

var (
	//go:embed order_details.html
	orderDetailsHTML string

	//go:embed order_details.js
	orderDetailsJS string
)

// escapeJSString escapes a string for safe embedding inside a JavaScript
// single-quoted string literal. Escapes backslash, single quote, and
// newlines to prevent JS injection (fixes #12).
func escapeJSString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}

func (u *ui) renderPage(r *http.Request) string {
	if u.Store() == nil {
		return shared.ToFlashError(u.CacheStore(), nil, r, "Shop store is not initialized", shared.AdminHomeURL(r), 10)
	}

	orderID := req.GetStringTrimmed(r, "order_id")
	if orderID == "" {
		return shared.ToFlashError(u.CacheStore(), nil, r, "Order ID is required", shared.URLR(r, shared.CONTROLLER_ORDERS, nil), 10)
	}

	// Escape user input for safe embedding (fixes #12 — XSS risk)
	orderIDHTML := html.EscapeString(orderID)
	orderIDJS := escapeJSString(orderID)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Orders", URL: shared.URLR(r, shared.CONTROLLER_ORDERS, nil)},
		{Name: "Order Details", URL: shared.URLR(r, shared.CONTROLLER_ORDER_DETAILS, map[string]string{"order_id": orderID})},
	})

	heading := hb.Heading1().HTML("Order Details: " + orderIDHTML)

	linksHelper := shared.NewLinksFromRequest(r)
	urlOrders := linksHelper.Orders(map[string]string{})
	urlLoadOrderDetails := linksHelper.OrderDetails(map[string]string{"action": "load-order-details", "order_id": orderID})

	// Escape URLs for JS string context too
	urlOrdersJS := escapeJSString(urlOrders)
	urlLoadOrderDetailsJS := escapeJSString(urlLoadOrderDetails)

	htmlContent := strings.ReplaceAll(orderDetailsHTML, "urlOrders", "'"+urlOrdersJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlLoadOrderDetails", "'"+urlLoadOrderDetailsJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "ORDER_ID", orderIDHTML)

	jsContent := strings.ReplaceAll(orderDetailsJS, "urlOrders", "'"+urlOrdersJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlLoadOrderDetails", "'"+urlLoadOrderDetailsJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "ORDER_ID", "'"+orderIDJS+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))

	return u.Layout(nil, r, "Order Details | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
	})
}
