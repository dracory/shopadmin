package shopadmin

import (
	"net/http"

	"github.com/dracory/shopadmin/shared"
)

// buildControllerRoutes creates the controller dispatch map.
// This is called once at construction time (in New()) or per-request
// (in Routes()). Each controller package exposes a UI(config) factory
// that returns a handler function.
func buildControllerRoutes(uiConfig shared.UiConfig, fileManagerURL string) map[string]func(w http.ResponseWriter, r *http.Request) {
	return map[string]func(w http.ResponseWriter, r *http.Request){
		shared.CONTROLLER_HOME:            homeHandler(uiConfig, fileManagerURL),
		shared.CONTROLLER_PRODUCTS:        productManagerHandler(uiConfig),
		shared.CONTROLLER_PRODUCT_UPDATE:  productUpdateHandler(uiConfig, fileManagerURL),
		shared.CONTROLLER_CATEGORIES:      categoryManagerHandler(uiConfig),
		shared.CONTROLLER_CATEGORY_CREATE: categoryCreateHandler(uiConfig),
		shared.CONTROLLER_CATEGORY_UPDATE: categoryUpdateHandler(uiConfig),
		shared.CONTROLLER_DISCOUNTS:       discountManagerHandler(uiConfig),
		shared.CONTROLLER_ORDERS:          orderManagerHandler(uiConfig),
		shared.CONTROLLER_ORDER_DETAILS:   orderDetailsHandler(uiConfig),
	}
}
