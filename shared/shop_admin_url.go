package shared

import "net/http"

// ShopAdminURL returns the shop admin base URL from request context
func ShopAdminURL(r *http.Request) string {
	value, ok := r.Context().Value(KeyShopAdminURL).(string)
	if !ok {
		return ""
	}
	return value
}
