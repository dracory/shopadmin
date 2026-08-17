package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/product_update"
	"github.com/dracory/shopadmin/shared"
)

// productUpdateHandler creates a handler function for the product update controller
func productUpdateHandler(uiConfig shared.UiConfig, fileManagerURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		product_update.UI(uiConfig, fileManagerURL).ProductUpdate(w, r)
	}
}
