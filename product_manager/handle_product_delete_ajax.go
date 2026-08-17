package product_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleProductDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		ProductID string `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.ProductID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Product ID is required").ToString()))
		return ""
	}

	if err := shopStore.ProductDeleteByID(ctx, reqBody.ProductID); err != nil {
		u.Logger().Error("Failed to delete product", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to delete product").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Product deleted successfully").ToString()))
	return ""
}
