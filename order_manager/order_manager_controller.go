package order_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadOrdersAjax = "load-orders-ajax"
)

// UiInterface defines the order manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	OrderManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new order manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// OrderManager handles the order manager controller requests
func (u *ui) OrderManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the order manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadOrdersAjax:
		return u.handleOrdersLoadAjax(w, r)
	default:
		return u.renderPage(r)
	}
}
