package discount_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

//go:embed details.html details.js
var detailsEmbed embed.FS

type discountDetailsComponent struct {
	ui         *ui
	request    *http.Request
	discount   shopstore.DiscountInterface
	discountID string
}

func NewDiscountDetailsComponent(u *ui) *discountDetailsComponent {
	return &discountDetailsComponent{ui: u}
}

func (c *discountDetailsComponent) Mount(r *http.Request, discount shopstore.DiscountInterface, discountID string) {
	c.request = r
	c.discount = discount
	c.discountID = discountID
}

func (c *discountDetailsComponent) Handle(r *http.Request) error {
	return nil
}

func (c *discountDetailsComponent) Render() hb.TagInterface {
	htmlBytes, _ := detailsEmbed.ReadFile("details.html")
	jsBytes, _ := detailsEmbed.ReadFile("details.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	urlDetailsLoad := shared.NewLinksFromRequest(c.request).DiscountUpdate(map[string]string{
		"action":      actionLoadDetails,
		"discount_id": c.discountID,
	})

	urlDetailsSave := shared.NewLinksFromRequest(c.request).DiscountUpdate(map[string]string{
		"action":      actionSaveDetails,
		"discount_id": c.discountID,
	})

	initScript := `
		const discountId = "` + c.discountID + `";
		const urlDetailsLoad = "` + urlDetailsLoad + `";
		const urlDetailsSave = "` + urlDetailsSave + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}
