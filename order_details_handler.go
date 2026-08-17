package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/order_details"
	"github.com/dracory/shopadmin/shared"
)

// orderDetailsHandler creates a handler function for the order details controller
func orderDetailsHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		order_details.UI(uiConfig).OrderDetails(w, r)
	}
}
