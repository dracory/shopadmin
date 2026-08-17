package product_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleProductDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		BulkProductIDs []string `json:"bulk_product_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.BulkProductIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("No product IDs provided").ToString()))
		return ""
	}

	for _, productID := range reqBody.BulkProductIDs {
		if err := shopStore.ProductDeleteByID(ctx, productID); err != nil {
			u.Logger().Error("Failed to delete product", "error", err, "product_id", productID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Products deleted successfully").ToString()))
	return ""
}
