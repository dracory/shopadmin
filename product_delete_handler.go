package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/product_delete"
	"github.com/dracory/shopadmin/shared"
)

// productDeleteHandler creates a handler function for the product delete controller
func productDeleteHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		product_delete.UI(uiConfig).ProductDelete(w, r)
	}
}
