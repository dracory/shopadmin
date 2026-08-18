package discount_manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

func (u *ui) handleDiscountCreateAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Method not allowed").ToString()))
		return ""
	}

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not configured").ToString()))
		return ""
	}

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	title := strings.TrimSpace(reqBody.Title)
	if title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Title is required").ToString()))
		return ""
	}

	discount := shopstore.NewDiscount()
	discount.SetTitle(title)

	if err := shopStore.DiscountCreate(r.Context(), discount); err != nil {
		u.Logger().Error("handleDiscountCreateAjax", "error", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to create discount").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Discount created successfully", map[string]interface{}{
		FieldDiscountID: discount.GetID(),
	}).ToString()))
	return ""
}
