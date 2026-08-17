package product_manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/shopstore"
)

func (u *ui) handleProductCreateAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	if u.Store() == nil {
		api.Respond(w, r, api.Error("Shop store not configured"))
		return ""
	}

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if strings.TrimSpace(reqBody.Title) == "" {
		api.Respond(w, r, api.Error("Title is required"))
		return ""
	}

	product := shopstore.NewProduct()
	product.SetTitle(strings.TrimSpace(reqBody.Title))

	if err := u.Store().ProductCreate(r.Context(), product); err != nil {
		if u.Logger() != nil {
			u.Logger().Error("productManagerController.handleProductCreateAjax", "error", err.Error())
		}
		api.Respond(w, r, api.Error("Failed to create product"))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData("Product created successfully", map[string]interface{}{FieldProductID: product.GetID()}))
	return ""
}
