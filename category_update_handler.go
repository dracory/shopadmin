package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/category_update"
	"github.com/dracory/shopadmin/shared"
)

// categoryUpdateHandler creates a handler function for the category update controller
func categoryUpdateHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		category_update.UI(uiConfig).CategoryUpdate(w, r)
	}
}
