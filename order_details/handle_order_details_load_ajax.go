package order_details

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

func (u *ui) handleOrderDetailsLoadAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		api.Respond(w, r, api.Error("Shop store not available"))
		return ""
	}

	var reqBody struct {
		OrderID string `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if reqBody.OrderID == "" {
		api.Respond(w, r, api.Error("Order ID is required"))
		return ""
	}

	order, err := shopStore.OrderFindByID(ctx, reqBody.OrderID)
	if err != nil {
		u.Logger().Error("Failed to load order", "error", err)
		api.Respond(w, r, api.Error("Failed to load order"))
		return ""
	}

	if order == nil {
		api.Respond(w, r, api.Error("Order not found"))
		return ""
	}

	customerName := ""
	customerEmail := ""

	// Resolve customer via CustomerResolver (replaces userstore dependency)
	resolver := u.CustomerResolver()
	if resolver != nil && order.GetCustomerID() != "" {
		customerName, customerEmail = resolver.FindByID(ctx, order.GetCustomerID())
	}

	// Fetch order line items filtered by order ID at the DB level
	// (fixes #6 — previously loaded ALL line items then filtered in Go)
	lineItemQuery := shopstore.NewOrderLineItemQuery().SetOrderID(order.GetID())
	lineItems, err := shopStore.OrderLineItemList(ctx, lineItemQuery)
	if err != nil {
		u.Logger().Error("Failed to load order line items", "error", err)
	}

	items := []map[string]any{}
	for _, item := range lineItems {
		productName := item.GetProductID()
		product, err := shopStore.ProductFindByID(ctx, item.GetProductID())
		if err == nil && product != nil {
			productName = product.GetTitle()
		}

		// Calculate line total = price * quantity (fixes #17 —
		// previously "total" was set to unit price, not line total)
		lineTotal := item.GetPriceFloat() * float64(item.GetQuantityInt())

		items = append(items, map[string]any{
			"id":       item.GetID(),
			"name":     productName,
			"quantity": item.GetQuantity(),
			"price":    item.GetPrice(),
			"total":    lineTotal,
		})
	}

	// Parse shipping address from memo
	var shippingAddress map[string]string
	if order.GetMemo() != "" {
		if err := json.Unmarshal([]byte(order.GetMemo()), &shippingAddress); err != nil {
			// Log for debugging — shippingAddress will be nil, frontend
			// will show empty address (fixes #22)
			u.Logger().Error("Failed to parse shipping address from memo",
				"error", err, "order_id", order.GetID())
		}
	}

	orderData := map[string]any{
		"id":               order.GetID(),
		"status":           order.GetStatus(),
		"created_at":       order.GetCreatedAt(),
		"updated_at":       order.GetUpdatedAt(),
		"customer_id":      order.GetCustomerID(),
		"customer_name":    customerName,
		"customer_email":   customerEmail,
		"items":            items,
		"shipping_address": shippingAddress,
	}

	api.Respond(w, r, api.SuccessWithData("Order details loaded successfully", map[string]any{
		"order": orderData,
	}))
	return ""
}
