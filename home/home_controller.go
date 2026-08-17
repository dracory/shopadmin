package home

import (
	"embed"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed *.html
//go:embed *.js
var homeFiles embed.FS

const (
	actionLoadStats = "load-stats"
)

// UiInterface defines the home controller's UI interface
type UiInterface interface {
	shared.UiInterface
	Home(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new home controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// Home handles the home controller requests
func (u *ui) Home(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the home controller request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")

	switch action {
	case actionLoadStats:
		return u.handleLoadStats(w, r)
	default:
		return u.renderPage(r)
	}
}

func (u *ui) renderPage(r *http.Request) string {
	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
	})

	heading := hb.Heading1().HTML("Shop Dashboard")

	htmlContent, err := homeFiles.ReadFile("home.html")
	if err != nil {
		u.Logger().Error("Failed to read home HTML template", "error", err)
		return hb.Div().HTML("Error loading home component").ToHTML()
	}

	jsContent, err := homeFiles.ReadFile("home.js")
	if err != nil {
		u.Logger().Error("Failed to read home JavaScript file", "error", err)
		return hb.Div().HTML("Error loading home component").ToHTML()
	}

	linksHelper := shared.NewLinksFromRequest(r)
	initScript := hb.Script(`
		const urlLoadStats = '` + linksHelper.Home(map[string]string{"action": actionLoadStats}) + `';
		const urlProducts = '` + linksHelper.Products(map[string]string{}) + `';
		const urlCategories = '` + linksHelper.Categories(map[string]string{}) + `';
		const urlDiscounts = '` + linksHelper.Discounts(map[string]string{}) + `';
		const urlOrders = '` + linksHelper.Orders(map[string]string{}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueContainer)

	return u.Layout(nil, r, "Shop | Dashboard", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
	})
}

func (u *ui) handleLoadStats(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if u.Store() == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	productCount, err := u.Store().ProductCount(ctx, shopstore.NewProductQuery())
	if err != nil {
		u.Logger().Error("Failed to count products", "error", err)
		productCount = 0
	}

	categoryCount, err := u.Store().CategoryCount(ctx, shopstore.NewCategoryQuery())
	if err != nil {
		u.Logger().Error("Failed to count categories", "error", err)
		categoryCount = 0
	}

	orderCount, err := u.Store().OrderCount(ctx, shopstore.NewOrderQuery())
	if err != nil {
		u.Logger().Error("Failed to count orders", "error", err)
		orderCount = 0
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.SuccessWithData("Stats loaded successfully", map[string]any{
		"product_count":  productCount,
		"category_count": categoryCount,
		"order_count":    orderCount,
	}).ToString()))
	return ""
}
