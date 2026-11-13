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
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
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

	linksHelper := shared.NewLinksFromRequest(r)
	initScript := hb.Script(`
		const urlLoadCategory = '` + linksHelper.CategoryUpdate(map[string]string{"action": actionLoadCategory, "category_id": categoryID}) + `';
		const urlUpdateCategory = '` + linksHelper.CategoryUpdate(map[string]string{"action": actionUpdateCategory, "category_id": categoryID}) + `';
		const urlCategories = '` + linksHelper.Categories(map[string]string{}) + `';
		const categoryID = '` + categoryID + `';
	`)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Categories. Update Category")).
		Child(initScript).
		Child(categoryUpdateForm()).
		Child(hb.Script(categoryUpdateJS))

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
		_, _ = w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	categoryID := req.GetStringTrimmed(r, "category_id")
	if categoryID == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Category ID is required").ToString()))
		return ""
	}

	category, err := shopStore.CategoryFindByID(ctx, categoryID)
	if err != nil || category == nil {
		u.Logger().Error("Failed to load category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Category not found").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Category loaded successfully", map[string]any{
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
		_, _ = w.Write([]byte(api.Error("Shop store not available").ToString()))
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
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.CategoryID == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Category ID is required").ToString()))
		return ""
	}

	category, err := shopStore.CategoryFindByID(ctx, reqBody.CategoryID)
	if err != nil || category == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Category not found").ToString()))
		return ""
	}

	category.SetTitle(reqBody.Title)
	category.SetDescription(reqBody.Description)
	category.SetStatus(reqBody.Status)
	category.SetParentID(reqBody.ParentID)

	if err := shopStore.CategoryUpdate(ctx, category); err != nil {
		u.Logger().Error("Failed to update category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to update category").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Category updated successfully").ToString()))
	return ""
}

// categoryUpdateForm renders the update-category form. Values are
// prefilled by categoryUpdateJS via the load-category action endpoint.
func categoryUpdateForm() hb.TagInterface {
	titleInput := hb.Input().
		Type("text").
		ID("category-title").
		Class("form-control").
		Attr("name", "title").
		Attr("required", "required")

	statusSelect := hb.Select().
		ID("category-status").
		Class("form-select").
		Attr("name", "status").
		Child(hb.Option().Value("active").Text("Active")).
		Child(hb.Option().Value("inactive").Text("Inactive")).
		Child(hb.Option().Value("draft").Text("Draft"))

	descriptionInput := hb.Textarea().
		ID("category-description").
		Class("form-control").
		Attr("name", "description").
		Attr("rows", "4")

	parentInput := hb.Input().
		Type("text").
		ID("category-parent-id").
		Class("form-control").
		Attr("name", "parent_id")

	submitButton := hb.Button().
		Type("button").
		ID("btn-update-category").
		Class("btn btn-primary").
		Attr("onclick", "submitCategoryUpdate()").
		HTML(`<i class="bi bi-check-lg me-1"></i> Save Category`)

	return hb.Div().
		Class("card bg-body-tertiary").
		Style("max-width: 640px;").
		Child(hb.Div().Class("card-body").
			Child(hb.Div().Class("mb-3").
				Child(hb.Label().Class("form-label").Attr("for", "category-title").Text("Title")).
				Child(titleInput)).
			Child(hb.Div().Class("mb-3").
				Child(hb.Label().Class("form-label").Attr("for", "category-status").Text("Status")).
				Child(statusSelect)).
			Child(hb.Div().Class("mb-3").
				Child(hb.Label().Class("form-label").Attr("for", "category-description").Text("Description")).
				Child(descriptionInput)).
			Child(hb.Div().Class("mb-3").
				Child(hb.Label().Class("form-label").Attr("for", "category-parent-id").Text("Parent Category")).
				Child(parentInput)).
			Child(submitButton))
}

// categoryUpdateJS loads the category into the form on page load and
// posts updates to the update-category endpoint, redirecting to the
// category manager on success.
const categoryUpdateJS = `
async function loadCategoryForm() {
	try {
		const response = await fetch(urlLoadCategory, {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({category_id: categoryID}),
		});
		const result = await response.json();
		if (result.status !== 'success' || !result.data || !result.data.category) {
			Swal.fire('Error', 'Failed to load category', 'error');
			return;
		}
		const c = result.data.category;
		document.getElementById('category-title').value = c.title || '';
		document.getElementById('category-status').value = c.status || 'draft';
		document.getElementById('category-description').value = c.description || '';
		document.getElementById('category-parent-id').value = c.parent_id || '';
	} catch (error) {
		Swal.fire('Error', 'Failed to load category', 'error');
	}
}

async function submitCategoryUpdate() {
	const title = document.getElementById('category-title').value.trim();
	if (!title) {
		Swal.fire('Validation', 'Title is required', 'warning');
		return;
	}
	const body = {
		category_id: categoryID,
		title: title,
		status: document.getElementById('category-status').value,
		description: document.getElementById('category-description').value,
		parent_id: document.getElementById('category-parent-id').value.trim(),
	};
	try {
		const response = await fetch(urlUpdateCategory, {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(body),
		});
		const result = await response.json();
		if (result.status === 'success') {
			window.location.href = urlCategories;
		} else {
			Swal.fire('Error', result.message || 'Failed to update category', 'error');
		}
	} catch (error) {
		Swal.fire('Error', 'Failed to update category', 'error');
	}
}

document.addEventListener('DOMContentLoaded', loadCategoryForm);
`
