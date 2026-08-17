package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/home"
	"github.com/dracory/shopadmin/shared"
)

// homeHandler creates a handler function for the home controller
func homeHandler(uiConfig shared.UiConfig, fileManagerURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		home.UI(uiConfig).Home(w, r)
	}
}
