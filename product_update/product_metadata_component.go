package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed metadata.html metadata.js
var metadataEmbed embed.FS

type productMetadataComponent struct {
	controller *ui
	request    *http.Request
	product    shopstore.ProductInterface
	productID  string

	formMetas map[string]string

	formErrorMessage   string
	formSuccessMessage string
}

func NewProductMetadataComponent(controller *ui) *productMetadataComponent {
	return &productMetadataComponent{controller: controller}
}

func (c *productMetadataComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID

	if product != nil {
		metas, _ := product.GetMetas()
		c.formMetas = metas
	}
}

func (c *productMetadataComponent) Handle(r *http.Request) error {
	c.formSuccessMessage = "Metadata saved successfully"
	return nil
}

func (c *productMetadataComponent) Render() hb.TagInterface {
	htmlBytes, _ := metadataEmbed.ReadFile("metadata.html")
	jsBytes, _ := metadataEmbed.ReadFile("metadata.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlMetadataLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-metadata",
		"product_id": c.productID,
	})

	urlMetadataSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-metadata",
		"product_id": c.productID,
	})

	initScript := `
		const productId = "` + c.productID + `";
		const urlMetadataLoad = "` + urlMetadataLoad + `";
		const urlMetadataSave = "` + urlMetadataSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
