package category_create

import (
	"encoding/json"
	"net/http"
	"uuid"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
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
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
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
		return shared.ErrorAlert("Shop store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Categories", URL: shared.URLR(r, shared.CONTROLLER_CATEGORIES, nil)},
		{Name: "Create Category", URL: ""},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	initScript := hb.Script(`
		const urlCreateCategory = '` + linksHelper.CategoryCreate(map[string]string{"action": actionCreateCategory}) + `';
		const urlCategories = '` + linksHelper.Categories(map[string]string{}) + `';
	`)

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Categories. Create Category")).
		Child(initScript).
		Child(categoryCreateForm()).
		Child(hb.Script(categoryCreateJS))

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
		_, _ = w.Write([]byte(api.Error("Shop store not available").ToString()))
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
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Title is required").ToString()))
		return ""
	}

	category := shopstore.NewCategory()
	category.SetID(uuid.NewV7().String())
	category.SetTitle(reqBody.Title)
	category.SetDescription(reqBody.Description)
	category.SetStatus(reqBody.Status)
	category.SetParentID(reqBody.ParentID)

	if err := shopStore.CategoryCreate(ctx, category); err != nil {
		u.Logger().Error("Failed to create category", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to create category").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Category created successfully", map[string]any{
		"category_id": category.GetID(),
	}).ToString()))
	return ""
}

// categoryCreateForm renders the create-category form (title, status,
// description, parent ID). Submitted via categoryCreateJS to the
// create-category action endpoint.
func categoryCreateForm() hb.TagInterface {
	titleInput := hb.Input().
		Type("text").
		ID("category-title").
		Class("form-control").
		Attr("name", "title").
		Attr("required", "required").
		Attr("placeholder", "Category title")

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
		Attr("rows", "4").
		Attr("placeholder", "Category description")

	parentInput := hb.Input().
		Type("text").
		ID("category-parent-id").
		Class("form-control").
		Attr("name", "parent_id").
		Attr("placeholder", "Parent category ID (optional)")

	submitButton := hb.Button().
		Type("button").
		ID("btn-create-category").
		Class("btn btn-primary").
		Attr("onclick", "submitCategoryCreate()").
		HTML(`<i class="bi bi-check-lg me-1"></i> Create Category`)

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

// categoryCreateJS posts the form to the create-category endpoint and
// redirects to the category manager on success.
const categoryCreateJS = `
async function submitCategoryCreate() {
	const title = document.getElementById('category-title').value.trim();
	if (!title) {
		Swal.fire('Validation', 'Title is required', 'warning');
		return;
	}
	const body = {
		title: title,
		status: document.getElementById('category-status').value,
		description: document.getElementById('category-description').value,
		parent_id: document.getElementById('category-parent-id').value.trim(),
	};
	try {
		const response = await fetch(urlCreateCategory, {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(body),
		});
		const result = await response.json();
		if (result.status === 'success') {
			window.location.href = urlCategories;
		} else {
			Swal.fire('Error', result.message || 'Failed to create category', 'error');
		}
	} catch (error) {
		Swal.fire('Error', 'Failed to create category', 'error');
	}
}
`
