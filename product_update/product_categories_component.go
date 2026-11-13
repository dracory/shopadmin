package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed categories.html categories.js
var categoriesEmbed embed.FS

// productCategoriesComponent renders the standalone "Categories" card
// shown below the product details card. It has its own load/save
// handlers and does not modify any other product data.
type productCategoriesComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string
}

// NewProductCategoriesComponent creates the categories card component.
func NewProductCategoriesComponent(controller *ui) *productCategoriesComponent {
	return &productCategoriesComponent{controller: controller}
}

// Mount binds the component to the request/product context.
func (c *productCategoriesComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

// Handle is a no-op; saving is performed via the save-categories AJAX action.
func (c *productCategoriesComponent) Handle(r *http.Request) error {
	return nil
}

// Render returns the categories card markup and its scripts.
func (c *productCategoriesComponent) Render() hb.TagInterface {
	htmlBytes, _ := categoriesEmbed.ReadFile("categories.html")
	jsBytes, _ := categoriesEmbed.ReadFile("categories.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlCategoriesLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-categories",
		"product_id": c.productID,
	})

	urlCategoriesSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-categories",
		"product_id": c.productID,
	})

	// Globals are prefixed with "categories" because sibling cards
	// (details, tags, media) also declare "const productId" and a
	// duplicate top-level const would fail the entire script block.
	initScript := `
		const categoriesProductId = "` + c.productID + `";
		const categoriesUrlLoad = "` + urlCategoriesLoad + `";
		const categoriesUrlSave = "` + urlCategoriesSave + `";
	`

	return hb.Div().
		ID("categories-app").
		Class("card mt-3").
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
