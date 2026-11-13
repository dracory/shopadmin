package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed admin_notes.html admin_notes.js
var adminNotesEmbed embed.FS

// productAdminNotesComponent renders the standalone "Admin Notes"
// card shown as the last card on the details view. It has its own
// load/save handlers and does not modify any other product data.
type productAdminNotesComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string
}

// NewProductAdminNotesComponent creates the admin notes card component.
func NewProductAdminNotesComponent(controller *ui) *productAdminNotesComponent {
	return &productAdminNotesComponent{controller: controller}
}

// Mount binds the component to the request/product context.
func (c *productAdminNotesComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

// Handle is a no-op; saving is performed via the save-notes AJAX action.
func (c *productAdminNotesComponent) Handle(r *http.Request) error {
	return nil
}

// Render returns the admin notes card markup and its scripts.
func (c *productAdminNotesComponent) Render() hb.TagInterface {
	htmlBytes, _ := adminNotesEmbed.ReadFile("admin_notes.html")
	jsBytes, _ := adminNotesEmbed.ReadFile("admin_notes.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlNotesLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-notes",
		"product_id": c.productID,
	})

	urlNotesSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-notes",
		"product_id": c.productID,
	})

	// Globals are prefixed with "adminNotes" because sibling cards
	// (details, tags, media) also declare "const productId" and a
	// duplicate top-level const would fail the entire script block.
	initScript := `
		const adminNotesProductId = "` + c.productID + `";
		const adminNotesUrlLoad = "` + urlNotesLoad + `";
		const adminNotesUrlSave = "` + urlNotesSave + `";
	`

	return hb.Div().
		ID("admin-notes-app").
		Class("card mt-3").
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
