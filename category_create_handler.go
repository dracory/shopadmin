package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/category_create"
	"github.com/dracory/shopadmin/shared"
)

// categoryCreateHandler creates a handler function for the category create controller
func categoryCreateHandler(uiConfig shared.UiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		category_create.UI(uiConfig).CategoryCreate(w, r)
	}
}
