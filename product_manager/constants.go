package product_manager

import (
	"context"
	"strings"

	"github.com/dracory/shopstore"
)

const (
	FieldID        = "id"
	FieldTitle     = "title"
	FieldStatus    = "status"
	FieldPrice     = "price"
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
	FieldProducts  = "products"
	FieldTotal     = "total"
	FieldProductID = "product_id"
)

// Condition field constants for the multi-filter system.
const (
	CondFieldSearch   = "search"
	CondFieldStatus   = "status"
	CondFieldID       = "id"
	CondFieldCategory = "category_id"
	CondFieldDateFrom = "date_from"
	CondFieldDateTo   = "date_to"
)

// condition is a single filter condition sent from the frontend.
type condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// applyConditions maps filter conditions to the product query.
// Multiple conditions on the same field: last one wins for single-value
// fields. Date conditions accumulate.
func applyConditions(query shopstore.ProductQueryInterface, conds []condition) {
	for _, c := range conds {
		val := strings.TrimSpace(c.Value)
		if val == "" {
			continue
		}
		switch c.Field {
		case CondFieldSearch:
			query.SetTitleLike(val)
		case CondFieldStatus:
			query.SetStatus(val)
		case CondFieldID:
			query.SetID(val)
		case CondFieldDateFrom:
			query.SetCreatedAtGte(val + " 00:00:00")
		case CondFieldDateTo:
			query.SetCreatedAtLte(val + " 23:59:59")
		}
	}
}

// categoryProductIDs resolves category conditions to product IDs via the
// product-category link table. Multiple category conditions: last one wins.
func categoryProductIDs(ctx context.Context, store shopstore.StoreInterface, conds []condition) ([]string, bool, error) {
	categoryID := ""
	for _, c := range conds {
		if c.Field != CondFieldCategory {
			continue
		}
		if val := strings.TrimSpace(c.Value); val != "" {
			categoryID = val
		}
	}

	if categoryID == "" {
		return nil, false, nil
	}

	links, err := store.ProductCategoryList(ctx, shopstore.NewProductCategoryQuery().
		SetCategoryID(categoryID))
	if err != nil {
		return nil, true, err
	}

	ids := make([]string, 0, len(links))
	for _, link := range links {
		ids = append(ids, link.GetProductID())
	}
	return ids, true, nil
}
