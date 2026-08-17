package product_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadProducts          = "load-products"
	actionProductDelete         = "delete-product"
	actionProductDeleteSelected = "delete-selected"
	actionCreateProduct         = "create-product-ajax"
)

// UiInterface defines the product manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	ProductManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new product manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// ProductManager handles the product manager controller requests
func (u *ui) ProductManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the product manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadProducts:
		return u.handleLoadProducts(w, r)
	case actionProductDelete:
		return u.handleProductDelete(w, r)
	case actionProductDeleteSelected:
		return u.handleProductDeleteSelected(w, r)
	case actionCreateProduct:
		return u.handleProductCreateAjax(w, r)
	default:
		return u.renderPage(w, r)
	}
}
