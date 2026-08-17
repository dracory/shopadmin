package category_update

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadCategory   = "load-category"
	actionUpdateCategory = "update-category"
)

// UiInterface defines the category update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	CategoryUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new category update controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// CategoryUpdate handles the category update controller requests
func (u *ui) CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the category update request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadCategory:
		return u.handleLoadCategory(w, r)
	case actionUpdateCategory:
		return u.handleUpdateCategory(w, r)
	default:
		return u.renderPage(r)
	}
}

func (u *ui) renderPage(r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	categoryID := req.GetStringTrimmed(r, "category_id")
	if categoryID == "" {
		return shared.ErrorAlert("Category ID is required")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Categories", URL: shared.URLR(r, shared.CONTROLLER_CATEGORIES, nil)},
		{Name: "Update Category", URL: ""},
	})

	heading := hb.Heading1().HTML("Update Category")

	linksHelper := shared.NewLinksFromRequest(r)
	initScript := hb.Script(`
		const urlLoadCategory = '` + linksHelper.CategoryUpdate(map[string]string{"action": actionLoadCategory, "category_id": categoryID}) + `';
		const urlUpdateCategory = '` + linksHelper.CategoryUpdate(map[string]string{"action": actionUpdateCategory, "category_id": categoryID}) + `';
		const categoryID = '` + categoryID + `';
	`)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(initScript).
		Child(hb.Div().ID("app"))

	return u.Layout(nil, r, "Update Category | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Htmx_1_9_4(),
			cdn.Sweetalert2_10(),
		},
	})
}

func (u *ui) handleLoadCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	categoryID := req.GetStringTrimmed(r, "category_id")
	if categoryID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Category ID is required").ToString()))
		return ""
	}

	category, err := shopStore.CategoryFindByID(ctx, categoryID)
	if err != nil || category == nil {
		u.Logger().Error("Failed to load category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Category not found").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Category loaded successfully", map[string]any{
		"category": map[string]any{
			"id":          category.GetID(),
			"title":       category.GetTitle(),
			"description": category.GetDescription(),
			"status":      category.GetStatus(),
			"parent_id":   category.GetParentID(),
		},
	}).ToString()))
	return ""
}

func (u *ui) handleUpdateCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		CategoryID  string `json:"category_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		ParentID    string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.CategoryID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Category ID is required").ToString()))
		return ""
	}

	category, err := shopStore.CategoryFindByID(ctx, reqBody.CategoryID)
	if err != nil || category == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Category not found").ToString()))
		return ""
	}

	category.SetTitle(reqBody.Title)
	category.SetDescription(reqBody.Description)
	category.SetStatus(reqBody.Status)
	category.SetParentID(reqBody.ParentID)

	if err := shopStore.CategoryUpdate(ctx, category); err != nil {
		u.Logger().Error("Failed to update category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to update category").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Category updated successfully").ToString()))
	return ""
}
