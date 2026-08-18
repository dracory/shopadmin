package discount_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadDiscounts          = "load-discounts"
	actionDiscountDelete         = "delete-discount"
	actionDiscountDeleteSelected = "delete-selected"
	actionCreateDiscount         = "create-discount-ajax"
)

// UiInterface defines the discount manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	DiscountManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new discount manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// DiscountManager handles the discount manager controller requests
func (u *ui) DiscountManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the discount manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadDiscounts:
		return u.handleLoadDiscounts(w, r)
	case actionDiscountDelete:
		return u.handleDiscountDelete(w, r)
	case actionDiscountDeleteSelected:
		return u.handleDiscountDeleteSelected(w, r)
	case actionCreateDiscount:
		return u.handleDiscountCreateAjax(w, r)
	default:
		return u.renderPage(w, r)
	}
}
