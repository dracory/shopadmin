package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/discount_manager"
	"github.com/dracory/shopadmin/shared"
)

// discountManagerHandler creates a handler function for the discount manager controller
func discountManagerHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		discount_manager.UI(uiConfig).DiscountManager(w, r)
	}
}
