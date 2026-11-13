package product_manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/neat"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
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

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadProducts":          linksHelper.Products(map[string]string{"action": actionLoadProducts}),
		"productDelete":         linksHelper.Products(map[string]string{"action": actionProductDelete}),
		"productDeleteSelected": linksHelper.Products(map[string]string{"action": actionProductDeleteSelected}),
		"viewProduct":           linksHelper.ProductView(map[string]string{}),
		"updateProduct":         linksHelper.ProductUpdate(map[string]string{}),
		"createProduct":         linksHelper.Products(map[string]string{"action": actionCreateProduct}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.productManagerUrls = %s;", string(urlsJSON)))

	categories, err := u.Store().CategoryList(r.Context(), shopstore.NewCategoryQuery().
		SetStatus(shopstore.CATEGORY_STATUS_ACTIVE).
		SetOrderBy(shopstore.COLUMN_TITLE).
		SetSortDirection(neat.SortAsc))
	categoryItems := []map[string]any{}
	if err != nil {
		u.Logger().Error("Failed to load categories for filter", "error", err)
	} else {
		for _, category := range categories {
			categoryItems = append(categoryItems, map[string]any{
				"id":    category.GetID(),
				"title": category.GetTitle(),
			})
		}
	}
	categoriesJSON, _ := json.Marshal(categoryItems)
	categoriesScript := hb.Script(fmt.Sprintf("window.productManagerCategories = %s;", string(categoriesJSON)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Products. Product Manager")).
		Child(urlsScript).
		Child(categoriesScript).
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
