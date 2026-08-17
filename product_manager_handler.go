package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/product_manager"
	"github.com/dracory/shopadmin/shared"
)

// productManagerHandler creates a handler function for the product manager controller
func productManagerHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		product_manager.UI(uiConfig).ProductManager(w, r)
	}
}
