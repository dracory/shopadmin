package order_details

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadOrderDetailsAjax = "load-order-details"
)

// UiInterface defines the order details controller's UI interface
type UiInterface interface {
	shared.UiInterface
	OrderDetails(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new order details controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// OrderDetails handles the order details controller requests
func (u *ui) OrderDetails(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the order details request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadOrderDetailsAjax:
		return u.handleOrderDetailsLoadAjax(w, r)
	default:
		return u.renderPage(r)
	}
}
