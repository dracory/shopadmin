package shared

import "net/http"

// Links provides URL helpers for shopadmin controllers.
// The base URL is read from request context (injected by Handle()),
// not hardcoded. This fixes the pre-existing issue where "/admin/shop"
// was hardcoded in every subcontroller.
type Links struct {
	baseURL string
}

// NewLinks creates a Links helper with the given base URL.
// If baseURL is empty, defaults to "/admin/shop".
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/shop"
	}
	return &Links{baseURL: baseURL}
}

// NewLinksFromRequest creates a Links helper using the shop admin URL
// from the request context.
func NewLinksFromRequest(r *http.Request) *Links {
	return NewLinks(ShopAdminURL(r))
}

// Home builds the URL for the home controller
func (l *Links) Home(params map[string]string) string {
	return l.url(CONTROLLER_HOME, params)
}

// Products builds the URL for the products controller
func (l *Links) Products(params map[string]string) string {
	return l.url(CONTROLLER_PRODUCTS, params)
}

// ProductUpdate builds the URL for the product update controller
func (l *Links) ProductUpdate(params map[string]string) string {
	return l.url(CONTROLLER_PRODUCT_UPDATE, params)
}

// Categories builds the URL for the categories controller
func (l *Links) Categories(params map[string]string) string {
	return l.url(CONTROLLER_CATEGORIES, params)
}

// CategoryCreate builds the URL for the category create controller
func (l *Links) CategoryCreate(params map[string]string) string {
	return l.url(CONTROLLER_CATEGORY_CREATE, params)
}

// CategoryUpdate builds the URL for the category update controller
func (l *Links) CategoryUpdate(params map[string]string) string {
	return l.url(CONTROLLER_CATEGORY_UPDATE, params)
}

// Discounts builds the URL for the discounts controller
func (l *Links) Discounts(params map[string]string) string {
	return l.url(CONTROLLER_DISCOUNTS, params)
}

// Orders builds the URL for the orders controller
func (l *Links) Orders(params map[string]string) string {
	return l.url(CONTROLLER_ORDERS, params)
}

// OrderDetails builds the URL for the order details controller
func (l *Links) OrderDetails(params map[string]string) string {
	return l.url(CONTROLLER_ORDER_DETAILS, params)
}

// url builds a URL for the given controller. The params map is copied
// before mutation (fixes pre-existing bug #10 where the caller's map
// was mutated).
func (l *Links) url(controller string, params map[string]string) string {
	return URL(l.baseURL, controller, params)
}
