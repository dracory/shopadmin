package shared

// Context keys for config values injected by Handle()
const KeyEndpoint = "endpoint"
const KeyAdminHomeURL = "admin_home_url"
const KeyShopAdminURL = "shop_admin_url"
const KeyFileManagerURL = "file_manager_url"

// Controller names used in the ?controller= query parameter
const (
	CONTROLLER_HOME            = "home"
	CONTROLLER_PRODUCTS        = "products"
	CONTROLLER_PRODUCT_VIEW    = "product_view"
	CONTROLLER_PRODUCT_UPDATE  = "product_update"
	CONTROLLER_CATEGORIES      = "categories"
	CONTROLLER_CATEGORY_CREATE = "category_create"
	CONTROLLER_CATEGORY_UPDATE = "category_update"
	CONTROLLER_DISCOUNTS       = "discounts"
	CONTROLLER_DISCOUNT_VIEW   = "discount_view"
	CONTROLLER_DISCOUNT_UPDATE = "discount_update"
	CONTROLLER_ORDERS          = "orders"
	CONTROLLER_ORDER_DETAILS   = "order_details"
)

// CatchAll is the catch-all route suffix
const CatchAll = "/*"

// Error messages
const ERROR_STORE_IS_NIL = "store cannot be nil"
const ERROR_LOGGER_IS_NIL = "logger cannot be nil"
