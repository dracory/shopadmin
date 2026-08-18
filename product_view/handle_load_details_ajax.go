package product_view

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

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Product ID is required").ToString()))
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to load product").ToString()))
		return ""
	}

	if product == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Product not found").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Product loaded successfully", map[string]any{
		"id":                product.GetID(),
		"title":             product.GetTitle(),
		"description":       product.GetDescription(),
		"short_description": product.GetShortDescription(),
		"price":             product.GetPrice(),
		"quantity":          product.GetQuantity(),
		"status":            product.GetStatus(),
		"memo":              product.GetMemo(),
		"parent_id":         product.GetParentID(),
		"created_at":        product.GetCreatedAt(),
		"updated_at":        product.GetUpdatedAt(),
		"price_float":       cast.ToFloat64(product.GetPrice()),
		"quantity_int":      product.GetQuantityInt(),
	}).ToString()))
	return ""
}
