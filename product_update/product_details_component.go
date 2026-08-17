package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed details.html details.js
var detailsEmbed embed.FS

type productDetailsComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string
}

func NewProductDetailsComponent(controller *ui) *productDetailsComponent {
	return &productDetailsComponent{controller: controller}
}

func (c *productDetailsComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

func (c *productDetailsComponent) Handle(r *http.Request) error {
	return nil
}

func (c *productDetailsComponent) Render() hb.TagInterface {
	htmlBytes, _ := detailsEmbed.ReadFile("details.html")
	jsBytes, _ := detailsEmbed.ReadFile("details.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlDetailsLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-details",
		"product_id": c.productID,
	})

	urlDetailsSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-details",
		"product_id": c.productID,
	})

	initScript := `
		const productId = "` + c.productID + `";
		const urlDetailsLoad = "` + urlDetailsLoad + `";
		const urlDetailsSave = "` + urlDetailsSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
