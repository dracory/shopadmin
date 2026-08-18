package order_manager

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
)

var (
	//go:embed orders.html
	ordersHTML string

	//go:embed orders.js
	ordersJS string
)

func (u *ui) renderPage(r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Orders", URL: shared.URLR(r, shared.CONTROLLER_ORDERS, nil)},
	})

	heading := hb.Heading1().HTML("Order Manager")

	linksHelper := shared.NewLinksFromRequest(r)
	urlLoadOrders := linksHelper.Orders(map[string]string{"action": actionLoadOrdersAjax})
	urlOrderDetails := linksHelper.OrderDetails(map[string]string{"order_id": "ORDER_ID"})

	html := strings.ReplaceAll(ordersHTML, "urlLoadOrders", "'"+urlLoadOrders+"'")
	html = strings.ReplaceAll(html, "urlOrderDetails", "'"+urlOrderDetails+"'")
	js := strings.ReplaceAll(ordersJS, "urlLoadOrders", "'"+urlLoadOrders+"'")
	js = strings.ReplaceAll(js, "urlOrderDetails", "'"+urlOrderDetails+"'")

	vueCDN := hb.Script("").Src(cdn.VueJs_3())

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	return u.Layout(nil, r, "Orders | Shop", content.ToHTML(), struct {
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
