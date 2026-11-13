package category_manager

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
	//go:embed categories.html
	categoriesHTML string

	//go:embed categories.js
	categoriesJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Categories", URL: shared.URLR(r, shared.CONTROLLER_CATEGORIES, nil)},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadCategories":         linksHelper.Categories(map[string]string{"action": actionLoadCategories}),
		"categoryDelete":         linksHelper.Categories(map[string]string{"action": actionCategoryDelete}),
		"categoryDeleteSelected": linksHelper.Categories(map[string]string{"action": actionCategoryDeleteSelected}),
		"categoryUpdate":         linksHelper.CategoryUpdate(map[string]string{}),
		"categoryCreate":         linksHelper.CategoryCreate(map[string]string{}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.categoryManagerUrls = %s;", string(urlsJSON)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Shop. Categories")).
		Child(urlsScript).
		Child(hb.Raw(categoriesHTML)).
		Child(hb.Script(categoriesJS))

	return u.Layout(w, r, "Categories | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Sweetalert2_10(),
		},
	})
}
