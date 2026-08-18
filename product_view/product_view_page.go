package product_view

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

	productID := req.GetStringTrimmed(r, "product_id")
	if productID == "" {
		return shared.ErrorAlert("Product ID is required")
	}

	// Escape user input for safe embedding (XSS protection)
	productIDHTML := html.EscapeString(productID)
	productIDJS := escapeJSString(productID)

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Products", URL: shared.URLR(r, shared.CONTROLLER_PRODUCTS, nil)},
		{Name: "View Product", URL: shared.URLR(r, shared.CONTROLLER_PRODUCT_VIEW, map[string]string{"product_id": productID})},
	})

	heading := hb.Heading1().HTML("View Product: " + productIDHTML)

	linksHelper := shared.NewLinksFromRequest(r)
	urlProducts := linksHelper.Products(map[string]string{})
	urlProductUpdate := linksHelper.ProductUpdate(map[string]string{"product_id": productID})
	urlLoadDetails := linksHelper.ProductView(map[string]string{"action": actionLoadDetails, "product_id": productID})
	urlProductDelete := linksHelper.Products(map[string]string{"action": "delete-product"})

	// Escape URLs for JS string context
	urlProductsJS := escapeJSString(urlProducts)
	urlProductUpdateJS := escapeJSString(urlProductUpdate)
	urlLoadDetailsJS := escapeJSString(urlLoadDetails)
	urlProductDeleteJS := escapeJSString(urlProductDelete)

	htmlContent := strings.ReplaceAll(viewHTML, "urlProducts", "'"+urlProductsJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlProductUpdate", "'"+urlProductUpdateJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlLoadDetails", "'"+urlLoadDetailsJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "urlProductDelete", "'"+urlProductDeleteJS+"'")
	htmlContent = strings.ReplaceAll(htmlContent, "PRODUCT_ID", productIDHTML)

	jsContent := strings.ReplaceAll(viewJS, "urlProducts", "'"+urlProductsJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlProductUpdate", "'"+urlProductUpdateJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlLoadDetails", "'"+urlLoadDetailsJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "urlProductDelete", "'"+urlProductDeleteJS+"'")
	jsContent = strings.ReplaceAll(jsContent, "PRODUCT_ID", "'"+productIDJS+"'")

	vueCDN := hb.Script("").Src(cdn.VueJs_3())

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))

	return u.Layout(nil, r, "View Product | Shop", content.ToHTML(), struct {
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
