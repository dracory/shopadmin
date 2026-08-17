package category_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadCategories         = "load-categories"
	actionCategoryDelete         = "delete-category"
	actionCategoryDeleteSelected = "delete-selected"
)

// UiInterface defines the category manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	CategoryManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new category manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// CategoryManager handles the category manager controller requests
func (u *ui) CategoryManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the category manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadCategories:
		return u.handleLoadCategories(w, r)
	case actionCategoryDelete:
		return u.handleCategoryDelete(w, r)
	case actionCategoryDeleteSelected:
		return u.handleCategoryDeleteSelected(w, r)
	default:
		return u.renderPage(w, r)
	}
}
