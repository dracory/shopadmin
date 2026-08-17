package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/category_manager"
	"github.com/dracory/shopadmin/shared"
)

// categoryManagerHandler creates a handler function for the category manager controller
func categoryManagerHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		category_manager.UI(uiConfig).CategoryManager(w, r)
	}
}
