package product_manager

import (
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
