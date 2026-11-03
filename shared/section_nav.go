package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
	"github.com/spf13/cast"
)

// SectionNav renders the shop admin section navigation as pills, matching the
// nav-pills pattern used across the admin packages (useradmin, blogadmin).
// The active tab is derived from the ?controller= query parameter.
// Orders, Products, and Discounts tabs show record counts as badges when the
// store is available. The File Manager tab is only shown when a file manager
// URL has been configured.
func SectionNav(r *http.Request, store shopstore.StoreInterface, logger *slog.Logger) hb.TagInterface {
	linksHelper := NewLinksFromRequest(r)

	active := "dashboard"
	controller := ""
	if r != nil && r.URL != nil {
		controller = r.URL.Query().Get("controller")
	}
	switch controller {
	case CONTROLLER_ORDERS, CONTROLLER_ORDER_DETAILS:
		active = "orders"
	case CONTROLLER_PRODUCTS, CONTROLLER_PRODUCT_VIEW, CONTROLLER_PRODUCT_UPDATE:
		active = "products"
	case CONTROLLER_CATEGORIES, CONTROLLER_CATEGORY_CREATE, CONTROLLER_CATEGORY_UPDATE:
		active = "categories"
	case CONTROLLER_DISCOUNTS, CONTROLLER_DISCOUNT_VIEW, CONTROLLER_DISCOUNT_UPDATE:
		active = "discounts"
	}

	tab := func(label, key string, url string, count int) hb.TagInterface {
		class := "nav-link"
		if key == active {
			class += " active"
		}
		link := hb.Hyperlink().HTML(label).Href(url).Class(class)
		if count >= 0 {
			link.Child(hb.Span().
				Class("badge bg-secondary ms-2").
				HTML(cast.ToString(count)))
		}
		return hb.LI().Class("nav-item").Child(link)
	}

	count := func(query func() (int64, error)) int {
		if store == nil {
			return -1
		}
		n, err := query()
		if err != nil {
			if logger != nil {
				logger.Error(err.Error())
			}
			return -1
		}
		return int(n)
	}

	ctx := r.Context()
	tabs := []hb.TagInterface{
		tab("Dashboard", "dashboard", linksHelper.Home(nil), -1),
		tab("Orders", "orders", linksHelper.Orders(nil),
			count(func() (int64, error) { return store.OrderCount(ctx, shopstore.NewOrderQuery()) })),
		tab("Products", "products", linksHelper.Products(nil),
			count(func() (int64, error) { return store.ProductCount(ctx, shopstore.NewProductQuery()) })),
		tab("Categories", "categories", linksHelper.Categories(nil), -1),
		tab("Discounts", "discounts", linksHelper.Discounts(nil),
			count(func() (int64, error) { return store.DiscountCount(ctx, shopstore.NewDiscountQuery()) })),
	}

	if fileManagerURL := FileManagerURL(r); fileManagerURL != "" {
		tabs = append(tabs, tab("File Manager", "files", fileManagerURL, -1))
	}

	ulNav := hb.UL().
		Class("nav nav-pills p-1 justify-content-center flex-grow-1").
		Children(tabs)

	label := hb.Span().
		Class("fw-bolder text-uppercase ms-3 me-3").
		Text("Shop").
		Child(hb.I().Class("bi bi-chevron-right ms-2"))

	return hb.Div().
		Class("card card-default mt-3 mb-4 bg-primary-subtle").
		Child(hb.Div().
			Class("card-body d-flex align-items-center").
			Style("padding: 2px;").
			Child(label).
			Child(ulNav))
}
