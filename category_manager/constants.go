package category_manager

import (
	"strings"

	"github.com/dracory/shopstore"
)

const (
	FieldID          = "id"
	FieldTitle       = "title"
	FieldDescription = "description"
	FieldStatus      = "status"
	FieldParentID    = "parent_id"
	FieldCreatedAt   = "created_at"
	FieldUpdatedAt   = "updated_at"
	FieldCategories  = "categories"
	FieldTotal       = "total"
	FieldCategoryID  = "category_id"
)

// Condition field constants for the multi-filter system.
const (
	CondFieldSearch = "search"
	CondFieldStatus = "status"
	CondFieldID     = "id"
)

// condition is a single filter condition sent from the frontend.
type condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// applyConditions maps filter conditions to the category query.
// Multiple conditions on the same field: last one wins.
func applyConditions(query shopstore.CategoryQueryInterface, conds []condition) {
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
		}
	}
}
