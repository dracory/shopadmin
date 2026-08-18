package discount_manager

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
	//go:embed discounts.html
	discountsHTML string

	//go:embed discounts.js
	discountsJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Discounts", URL: shared.URLR(r, shared.CONTROLLER_DISCOUNTS, nil)},
	})

	heading := hb.Heading1().HTML("Discount Manager")

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadDiscounts":  linksHelper.Discounts(map[string]string{"action": actionLoadDiscounts}),
		"createDiscount": linksHelper.Discounts(map[string]string{"action": actionCreateDiscount}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.discountManagerUrls = %s;", string(urlsJSON)))

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(urlsScript).
		Child(hb.Raw(discountsHTML)).
		Child(hb.Script(discountsJS))

	return u.Layout(w, r, "Discounts | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Sweetalert2_10(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
