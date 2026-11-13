package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed variants.html variants.js
var variantsEmbed embed.FS

// productVariantsComponent renders the standalone "Variants" card
// shown between the Categories and Admin Notes cards. It manages the
// parent's variant dimension schema and one child product per
// dimension combination through its own load/save handlers.
type productVariantsComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string
}

// NewProductVariantsComponent creates the variants card component.
func NewProductVariantsComponent(controller *ui) *productVariantsComponent {
	return &productVariantsComponent{controller: controller}
}

// Mount binds the component to the request/product context.
func (c *productVariantsComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

// Handle is a no-op; saving is performed via the save-variants AJAX action.
func (c *productVariantsComponent) Handle(r *http.Request) error {
	return nil
}

// Render returns the variants card markup and its scripts.
func (c *productVariantsComponent) Render() hb.TagInterface {
	htmlBytes, _ := variantsEmbed.ReadFile("variants.html")
	jsBytes, _ := variantsEmbed.ReadFile("variants.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlVariantsLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-variants",
		"product_id": c.productID,
	})

	urlVariantsSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-variants",
		"product_id": c.productID,
	})

	// Globals are prefixed with "variants" because sibling cards
	// (details, tags, media) also declare "const productId" and a
	// duplicate top-level const would fail the entire script block.
	initScript := `
		const variantsProductId = "` + c.productID + `";
		const variantsUrlLoad = "` + urlVariantsLoad + `";
		const variantsUrlSave = "` + urlVariantsSave + `";
	`

	return hb.Div().
		ID("variants-app").
		Class("card mt-3").
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
