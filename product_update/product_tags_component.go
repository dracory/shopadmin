package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed tags.html tags.js
var tagsEmbed embed.FS

type productTagsComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string
}

func NewProductTagsComponent(controller *ui) *productTagsComponent {
	return &productTagsComponent{controller: controller}
}

func (c *productTagsComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
}

func (c *productTagsComponent) Handle(r *http.Request) error {
	return nil
}

func (c *productTagsComponent) Render() hb.TagInterface {
	htmlBytes, _ := tagsEmbed.ReadFile("tags.html")
	jsBytes, _ := tagsEmbed.ReadFile("tags.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlTagsLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-tags",
		"product_id": c.productID,
	})

	urlTagsSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-tags",
		"product_id": c.productID,
	})

	initScript := `
		const productId = "` + c.productID + `";
		const urlTagsLoad = "` + urlTagsLoad + `";
		const urlTagsSave = "` + urlTagsSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
