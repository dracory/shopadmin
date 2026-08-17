package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopstore"
	"github.com/spf13/cast"
)

// Header renders the shop admin navigation header.
// It is nil-safe for both store and logger (fixes pre-existing bug #11).
func Header(store shopstore.StoreInterface, logger *slog.Logger, r *http.Request) hb.TagInterface {
	if store == nil {
		if logger != nil {
			logger.Error("shop store is nil")
		}
		return nil
	}

	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(AdminHomeURL(r)).
		Class("nav-link")

	linkShop := hb.NewHyperlink().
		HTML("Shop").
		Href(URLR(r, CONTROLLER_HOME, nil)).
		Class("nav-link")

	linkOrders := hb.Hyperlink().
		HTML("Orders").
		Href(URLR(r, CONTROLLER_ORDERS, nil)).
		Class("nav-link")

	linkDiscounts := hb.Hyperlink().
		HTML("Discounts").
		Href(URLR(r, CONTROLLER_DISCOUNTS, nil)).
		Class("nav-link")

	linkProducts := hb.Hyperlink().
		HTML("Products ").
		Href(URLR(r, CONTROLLER_PRODUCTS, nil)).
		Class("nav-link")

	productsCount, err := store.ProductCount(r.Context(), shopstore.NewProductQuery())
	if err != nil {
		if logger != nil {
			logger.Error(err.Error())
		}
		productsCount = -1
	}

	ordersCount, err := store.OrderCount(r.Context(), shopstore.NewOrderQuery())
	if err != nil {
		if logger != nil {
			logger.Error(err.Error())
		}
		ordersCount = -1
	}

	discountsCount, err := store.DiscountCount(r.Context(), shopstore.NewDiscountQuery())
	if err != nil {
		if logger != nil {
			logger.Error(err.Error())
		}
		discountsCount = -1
	}

	ulNav := hb.NewUL().
		Class("nav  nav-pills justify-content-center").
		Child(hb.NewLI().
			Class("nav-item").Child(linkHome)).
		Child(hb.NewLI().
			Class("nav-item").Child(linkShop)).
		Child(hb.LI().
			Class("nav-item").
			Child(linkOrders.
				Child(hb.Span().
					Class("badge bg-secondary ms-2").
					HTML(cast.ToString(ordersCount))))).
		Child(hb.LI().
			Class("nav-item").
			Child(linkProducts.
				Child(hb.Span().
					Class("badge bg-secondary ms-2").
					HTML(cast.ToString(productsCount))))).
		Child(hb.LI().
			Class("nav-item").
			Child(linkDiscounts.
				Child(hb.Span().
					Class("badge bg-secondary ms-2").
					HTML(cast.ToString(discountsCount)))))

	fileManagerURL := FileManagerURL(r)
	if fileManagerURL != "" {
		linkFileManager := hb.Hyperlink().
			HTML("File Manager").
			Href(fileManagerURL).
			Class("nav-link")
		ulNav.Child(hb.LI().
			Class("nav-item").
			Child(linkFileManager))
	}

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
}
