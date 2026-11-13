package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed status.html status.js
var statusEmbed embed.FS

// productStatusComponent renders the standalone "Status" card shown
// above the other cards. It has its own load/save handlers and does
// not modify any other product data.
type productStatusComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string
}

// NewProductStatusComponent creates the status card component.
func NewProductStatusComponent(controller *ui) *productStatusComponent {
	return &productStatusComponent{controller: controller}
}

// Mount binds the component to the request/product context.
func (c *productStatusComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

// Handle is a no-op; saving is performed via the save-status AJAX action.
func (c *productStatusComponent) Handle(r *http.Request) error {
	return nil
}

// Render returns the status card markup and its scripts.
func (c *productStatusComponent) Render() hb.TagInterface {
	htmlBytes, _ := statusEmbed.ReadFile("status.html")
	jsBytes, _ := statusEmbed.ReadFile("status.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlStatusLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-status",
		"product_id": c.productID,
	})

	urlStatusSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-status",
		"product_id": c.productID,
	})

	// Globals are prefixed with "status" because sibling cards
	// (details, tags, media) also declare "const productId" and a
	// duplicate top-level const would fail the entire script block.
	initScript := `
		const statusProductId = "` + c.productID + `";
		const statusUrlLoad = "` + urlStatusLoad + `";
		const statusUrlSave = "` + urlStatusSave + `";
	`

	return hb.Div().
		ID("status-app").
		Class("card mb-3").
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
