package category_manager

import (
	_ "embed"
	"net/http"
	"strings"

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
		return shared.ToFlashError(u.CacheStore(), w, r, "Shop store is not initialized", shared.AdminHomeURL(r), 10)
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Categories", URL: shared.URLR(r, shared.CONTROLLER_CATEGORIES, nil)},
	})

	heading := hb.Heading1().HTML("Category Manager")

	linksHelper := shared.NewLinksFromRequest(r)
	urlLoadCategories := linksHelper.Categories(map[string]string{"action": actionLoadCategories})
	urlCategoryDelete := linksHelper.Categories(map[string]string{"action": actionCategoryDelete})
	urlCategoryDeleteSelected := linksHelper.Categories(map[string]string{"action": actionCategoryDeleteSelected})

	html := strings.ReplaceAll(categoriesHTML, "urlLoadCategories", "'"+urlLoadCategories+"'")
	html = strings.ReplaceAll(html, "urlCategoryDelete", "'"+urlCategoryDelete+"'")
	html = strings.ReplaceAll(html, "urlCategoryDeleteSelected", "'"+urlCategoryDeleteSelected+"'")

	js := strings.ReplaceAll(categoriesJS, "urlLoadCategories", "'"+urlLoadCategories+"'")
	js = strings.ReplaceAll(js, "urlCategoryDelete", "'"+urlCategoryDelete+"'")
	js = strings.ReplaceAll(js, "urlCategoryDeleteSelected", "'"+urlCategoryDeleteSelected+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	return u.Layout(w, r, "Categories | Shop", content.ToHTML(), struct {
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
