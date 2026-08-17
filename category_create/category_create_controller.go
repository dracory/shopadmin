package category_create

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
	"github.com/dracory/uid"
)

const (
	actionCreateCategory = "create-category"
)

// UiInterface defines the category create controller's UI interface
type UiInterface interface {
	shared.UiInterface
	CategoryCreate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new category create controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// CategoryCreate handles the category create controller requests
func (u *ui) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the category create request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionCreateCategory:
		return u.handleCreateCategory(w, r)
	default:
		return u.renderPage(r)
	}
}

func (u *ui) renderPage(r *http.Request) string {
	if u.Store() == nil {
		return shared.ToFlashError(u.CacheStore(), nil, r, "Shop store is not initialized", shared.AdminHomeURL(r), 10)
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Categories", URL: shared.URLR(r, shared.CONTROLLER_CATEGORIES, nil)},
		{Name: "Create Category", URL: ""},
	})

	heading := hb.Heading1().HTML("Create Category")

	linksHelper := shared.NewLinksFromRequest(r)
	initScript := hb.Script(`
		const urlCreateCategory = '` + linksHelper.CategoryCreate(map[string]string{"action": actionCreateCategory}) + `';
	`)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(initScript).
		Child(hb.Div().ID("app"))

	return u.Layout(nil, r, "Create Category | Shop", content.ToHTML(), struct {
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

func (u *ui) handleCreateCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
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

	if reqBody.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Title is required").ToString()))
		return ""
	}

	category := shopstore.NewCategory()
	category.SetID(uid.HumanUid())
	category.SetTitle(reqBody.Title)
	category.SetDescription(reqBody.Description)
	category.SetStatus(reqBody.Status)
	category.SetParentID(reqBody.ParentID)

	if err := shopStore.CategoryCreate(ctx, category); err != nil {
		u.Logger().Error("Failed to create category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to create category").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Category created successfully", map[string]any{
		"category_id": category.GetID(),
	}).ToString()))
	return ""
}
