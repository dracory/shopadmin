package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/order_manager"
	"github.com/dracory/shopadmin/shared"
)

// orderManagerHandler creates a handler function for the order manager controller
func orderManagerHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		order_manager.UI(uiConfig).OrderManager(w, r)
	}
}
