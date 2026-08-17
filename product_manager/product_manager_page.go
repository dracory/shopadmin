package product_manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
)

var (
	//go:embed products.html
	productsHTML string

	//go:embed products.js
	productsJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Products", URL: shared.URLR(r, shared.CONTROLLER_PRODUCTS, nil)},
	})

	heading := hb.Heading1().HTML("Product Manager")

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadProducts":          linksHelper.Products(map[string]string{"action": actionLoadProducts}),
		"productDelete":         linksHelper.Products(map[string]string{"action": actionProductDelete}),
		"productDeleteSelected": linksHelper.Products(map[string]string{"action": actionProductDeleteSelected}),
		"updateProduct":         linksHelper.ProductUpdate(map[string]string{}),
		"createProduct":         linksHelper.Products(map[string]string{"action": actionCreateProduct}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.productManagerUrls = %s;", string(urlsJSON)))

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(urlsScript).
		Child(hb.Raw(productsHTML)).
		Child(hb.Script(productsJS))

	return u.Layout(w, r, "Products | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
