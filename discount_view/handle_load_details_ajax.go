package discount_view

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/spf13/cast"
)

func (u *ui) handleLoadDetails(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	discountID := r.URL.Query().Get("discount_id")
	if discountID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Discount ID is required").ToString()))
		return ""
	}

	discount, err := shopStore.DiscountFindByID(ctx, discountID)
	if err != nil {
		u.Logger().Error("Failed to load discount", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to load discount").ToString()))
		return ""
	}

	if discount == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Discount not found").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Discount loaded successfully", map[string]any{
		"id":                    discount.GetID(),
		"code":                  discount.GetCode(),
		"title":                 discount.GetTitle(),
		"description":           discount.GetDescription(),
		"type":                  discount.GetType(),
		"amount":                cast.ToFloat64(discount.GetAmount()),
		"status":                discount.GetStatus(),
		"starts_at":             discount.GetStartsAt(),
		"ends_at":               discount.GetEndsAt(),
		"memo":                  discount.GetMemo(),
		"created_at":            discount.GetCreatedAt(),
		"updated_at":            discount.GetUpdatedAt(),
		"max_uses":              discount.GetMaxUses(),
		"max_uses_count":        discount.GetMaxUsesCount(),
		"max_uses_per_customer": discount.GetMaxUsesPerCustomer(),
	}).ToString()))
	return ""
}
