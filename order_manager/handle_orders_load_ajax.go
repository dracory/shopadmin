package order_manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/neat"
	"github.com/dracory/shopstore"
)

func (u *ui) handleOrdersLoadAjax(w http.ResponseWriter, r *http.Request) string {
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
		Page          int    `json:"page"`
		PerPage       int    `json:"per_page"`
		SortBy        string `json:"sort_by"`
		Sort          string `json:"sort"`
		Status        string `json:"status"`
		CustomerName  string `json:"customer_name"`
		CustomerEmail string `json:"customer_email"`
		OrderID       string `json:"order_id"`
		CreatedFrom   string `json:"created_from"`
		CreatedTo     string `json:"created_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if reqBody.Page < 0 {
		reqBody.Page = 0
	}
	if reqBody.PerPage <= 0 {
		reqBody.PerPage = 10
	}
	if reqBody.SortBy == "" {
		reqBody.SortBy = shopstore.COLUMN_CREATED_AT
	}
	if reqBody.Sort == "" {
		reqBody.Sort = neat.SortDesc
	}

	// Load ALL orders without pagination, filter in memory, then slice
	// for the current page (fixes #15 — previously the query applied
	// Offset/Limit at the DB level, so len(filteredOrders) was at most
	// PerPage, making the total count useless for pagination)
	//
	// PERF TRADEOFF (#21): This loads all orders into memory on every
	// request. Acceptable for small-to-medium stores. The proper fix
	// is to move status/customer/date filters into the DB query layer
	// (shopstore.OrderQuery), then use OrderCount with the filtered
	// query and OrderList with Offset/Limit for pagination.
	query := shopstore.NewOrderQuery().
		SetOrderBy(reqBody.SortBy).
		SetSortDirection(reqBody.Sort)

	orders, err := shopStore.OrderList(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to load orders", "error", err)
		api.Respond(w, r, api.Error("Failed to load orders"))
		return ""
	}

	filteredOrders := orders

	// Filter by status
	if reqBody.Status != "" {
		filteredOrders = []shopstore.OrderInterface{}
		for _, order := range orders {
			if order.GetStatus() == reqBody.Status {
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	// Filter by order ID
	if reqBody.OrderID != "" {
		tempOrders := filteredOrders
		filteredOrders = []shopstore.OrderInterface{}
		for _, order := range tempOrders {
			if strings.Contains(strings.ToLower(order.GetID()), strings.ToLower(reqBody.OrderID)) {
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	// Filter by customer name/email via CustomerResolver (replaces userstore dependency)
	resolver := u.CustomerResolver()
	if (reqBody.CustomerName != "" || reqBody.CustomerEmail != "") && resolver != nil {
		matchingCustomerIDs, err := resolver.SearchIDs(ctx, reqBody.CustomerName, reqBody.CustomerEmail)
		if err != nil {
			u.Logger().Error("Failed to search customer IDs", "error", err)
			matchingCustomerIDs = nil
		}

		if len(matchingCustomerIDs) > 0 {
			tempOrders := filteredOrders
			filteredOrders = []shopstore.OrderInterface{}
			for _, order := range tempOrders {
				for _, customerID := range matchingCustomerIDs {
					if order.GetCustomerID() == customerID {
						filteredOrders = append(filteredOrders, order)
						break
					}
				}
			}
		} else {
			filteredOrders = []shopstore.OrderInterface{}
		}
	}

	// Filter by date range
	if reqBody.CreatedFrom != "" || reqBody.CreatedTo != "" {
		tempOrders := filteredOrders
		filteredOrders = []shopstore.OrderInterface{}
		for _, order := range tempOrders {
			createdAt := order.GetCreatedAt()
			if createdAt == "" {
				continue
			}

			match := true
			if reqBody.CreatedFrom != "" && createdAt < reqBody.CreatedFrom {
				match = false
			}
			if reqBody.CreatedTo != "" && createdAt > reqBody.CreatedTo {
				match = false
			}

			if match {
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	// Total is the count of ALL filtered orders across all pages
	total := len(filteredOrders)

	// Slice for the current page
	offset := reqBody.Page * reqBody.PerPage
	if offset > len(filteredOrders) {
		offset = len(filteredOrders)
	}
	end := offset + reqBody.PerPage
	if end > len(filteredOrders) {
		end = len(filteredOrders)
	}
	pagedOrders := filteredOrders[offset:end]

	orderList := []map[string]any{}
	for _, order := range pagedOrders {
		customerName := ""
		customerEmail := ""

		// Resolve customer via CustomerResolver (replaces userstore dependency)
		if resolver != nil && order.GetCustomerID() != "" {
			customerName, customerEmail = resolver.FindByID(ctx, order.GetCustomerID())
		}

		orderList = append(orderList, map[string]any{
			FieldID:            order.GetID(),
			FieldStatus:        order.GetStatus(),
			FieldCreatedAt:     order.GetCreatedAt(),
			FieldUpdatedAt:     order.GetUpdatedAt(),
			FieldCustomerID:    order.GetCustomerID(),
			FieldCustomerName:  customerName,
			FieldCustomerEmail: customerEmail,
		})
	}

	api.Respond(w, r, api.SuccessWithData("Orders loaded successfully", map[string]any{
		FieldOrders: orderList,
		FieldTotal:  total,
		"page":      reqBody.Page,
		"per_page":  reqBody.PerPage,
	}))
	return ""
}
