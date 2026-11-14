package shared

import "testing"

func TestControllerConstants(t *testing.T) {
	for name, c := range map[string]string{
		"home":            CONTROLLER_HOME,
		"products":        CONTROLLER_PRODUCTS,
		"product_view":    CONTROLLER_PRODUCT_VIEW,
		"product_update":  CONTROLLER_PRODUCT_UPDATE,
		"categories":      CONTROLLER_CATEGORIES,
		"category_create": CONTROLLER_CATEGORY_CREATE,
		"category_update": CONTROLLER_CATEGORY_UPDATE,
		"discounts":       CONTROLLER_DISCOUNTS,
		"discount_view":   CONTROLLER_DISCOUNT_VIEW,
		"discount_update": CONTROLLER_DISCOUNT_UPDATE,
		"orders":          CONTROLLER_ORDERS,
		"order_details":   CONTROLLER_ORDER_DETAILS,
		"settings":        CONTROLLER_SETTINGS,
	} {
		if c == "" || c != name {
			t.Fatalf("unexpected controller constant %q for %q", c, name)
		}
	}

	if CatchAll != "/*" {
		t.Fatalf("unexpected CatchAll %q", CatchAll)
	}
	if ERROR_STORE_IS_NIL == "" || ERROR_LOGGER_IS_NIL == "" {
		t.Fatal("expected non-empty error constants")
	}
}
