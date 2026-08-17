package discount_manager

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
)

var (
	//go:embed discounts.html
	discountsHTML string

	//go:embed discounts.js
	discountsJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ToFlashError(u.CacheStore(), w, r, "Shop store is not initialized", shared.AdminHomeURL(r), 10)
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Discounts", URL: shared.URLR(r, shared.CONTROLLER_DISCOUNTS, nil)},
	})

	heading := hb.Heading1().HTML("Discount Manager")

	linksHelper := shared.NewLinksFromRequest(r)
	urlLoadDiscounts := linksHelper.Discounts(map[string]string{"action": actionLoadDiscounts})
	urlDiscountDelete := linksHelper.Discounts(map[string]string{"action": actionDiscountDelete})
	urlDiscountDeleteSelected := linksHelper.Discounts(map[string]string{"action": actionDiscountDeleteSelected})

	html := strings.ReplaceAll(discountsHTML, "urlLoadDiscounts", "'"+urlLoadDiscounts+"'")
	html = strings.ReplaceAll(html, "urlDiscountDelete", "'"+urlDiscountDelete+"'")
	html = strings.ReplaceAll(html, "urlDiscountDeleteSelected", "'"+urlDiscountDeleteSelected+"'")

	js := strings.ReplaceAll(discountsJS, "urlLoadDiscounts", "'"+urlLoadDiscounts+"'")
	js = strings.ReplaceAll(js, "urlDiscountDelete", "'"+urlDiscountDelete+"'")
	js = strings.ReplaceAll(js, "urlDiscountDeleteSelected", "'"+urlDiscountDeleteSelected+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	return u.Layout(w, r, "Discounts | Shop", content.ToHTML(), struct {
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
