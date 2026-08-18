package discount_view

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
	//go:embed view.html
	viewHTML string

	//go:embed view.js
	viewJS string
)

// escapeJSString escapes a string for safe embedding inside a JavaScript
// single-quoted string literal. Escapes backslash, single quote, and
// newlines to prevent JS injection.
func escapeJSString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}

func (u *ui) renderPage(r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	discountID := req.GetStringTrimmed(r, "discount_id")
	if discountID == "" {
		return shared.ErrorAlert("Discount ID is required")
	}

	// Escape user input for safe embedding (XSS protection)
	discountIDHTML := html.EscapeString(discountID)
	discountIDJS := escapeJSString(discountID)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Discounts", URL: shared.URLR(r, shared.CONTROLLER_DISCOUNTS, nil)},
		{Name: "View Discount", URL: shared.URLR(r, shared.CONTROLLER_DISCOUNT_VIEW, map[string]string{"discount_id": discountID})},
	})

	heading := hb.Heading1().HTML("View Discount: " + discountIDHTML)

	linksHelper := shared.NewLinksFromRequest(r)
	urlDiscounts := linksHelper.Discounts(map[string]string{})
	urlDiscountUpdate := linksHelper.DiscountUpdate(map[string]string{"discount_id": discountID})
	urlLoadDetails := linksHelper.DiscountView(map[string]string{"action": actionLoadDetails, "discount_id": discountID})
	urlDiscountDelete := linksHelper.Discounts(map[string]string{"action": "delete-discount"})

	// Escape URLs for JS string context
	urlDiscountsJS := escapeJSString(urlDiscounts)
	urlDiscountUpdateJS := escapeJSString(urlDiscountUpdate)
	urlLoadDetailsJS := escapeJSString(urlLoadDetails)
	urlDiscountDeleteJS := escapeJSString(urlDiscountDelete)

	htmlContent := strings.ReplaceAll(viewHTML, "urlDiscounts", "'"+urlDiscountsJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlDiscountUpdate", "'"+urlDiscountUpdateJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlLoadDetails", "'"+urlLoadDetailsJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlDiscountDelete", "'"+urlDiscountDeleteJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "DISCOUNT_ID", discountIDHTML)

	jsContent := strings.ReplaceAll(viewJS, "urlDiscounts", "'"+urlDiscountsJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlDiscountUpdate", "'"+urlDiscountUpdateJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlLoadDetails", "'"+urlLoadDetailsJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlDiscountDelete", "'"+urlDiscountDeleteJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "DISCOUNT_ID", "'"+discountIDJS+"'")

	vueCDN := hb.Script("").Src(cdn.VueJs_3())

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))

	return u.Layout(nil, r, "View Discount | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
