package product_update

import (
	"embed"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
	"github.com/spf13/cast"
)

//go:embed media.html media.js
var mediaEmbed embed.FS

type productMediaComponent struct {
	controller     *ui
	request        *http.Request
	product        shopstore.ProductInterface
	productID      string
	fileManagerURL string

	formMedias []shopstore.MediaInterface

	formErrorMessage   string
	formSuccessMessage string
}

func NewProductMediaComponent(controller *ui) *productMediaComponent {
	return &productMediaComponent{controller: controller}
}

func (c *productMediaComponent) Mount(r *http.Request, product shopstore.ProductInterface, productID string) {
	c.request = r
	c.product = product
	c.productID = productID
	// Use the fileManagerURL from the controller, not a hardcoded value
	// (fixes #7)
	c.fileManagerURL = c.controller.fileManagerURL

	store := c.controller.Store()
	if store == nil {
		return
	}

	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	medias, _ := store.MediaList(r.Context(), mediaQuery)
	c.formMedias = medias
}

func (c *productMediaComponent) Handle(r *http.Request) error {
	store := c.controller.Store()
	if store == nil {
		c.formErrorMessage = "Shop store not available"
		return nil
	}

	ctx := r.Context()

	mediaURLs := req.GetAll(r)["media_url"]
	mediaTypes := req.GetAll(r)["media_type"]

	if len(mediaURLs) != len(mediaTypes) {
		c.formErrorMessage = "Media URLs and types count mismatch"
		return nil
	}

	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(c.productID)
	existingMedias, err := store.MediaList(ctx, mediaQuery)
	if err != nil {
		c.controller.Logger().Error("At productMediaComponent > Handle", "error", err.Error())
		c.formErrorMessage = "System error. Loading existing media failed"
		return nil
	}

	for _, existingMedia := range existingMedias {
		err := store.MediaDelete(ctx, existingMedia)
		if err != nil {
			c.controller.Logger().Error("At productMediaComponent > Handle", "error", err.Error())
		}
	}

	for i, mediaURL := range mediaURLs {
		if mediaURL == "" {
			continue
		}

		if len(mediaTypes) <= i || mediaTypes[i] == "" {
			continue
		}

		media := shopstore.NewMedia()
		media.SetEntityID(c.productID)
		media.SetURL(mediaURL)
		media.SetType(mediaTypes[i])
		media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
		media.SetSequence(i)

		err := store.MediaCreate(ctx, media)
		if err != nil {
			c.controller.Logger().Error("At productMediaComponent > Handle", "error", err.Error())
			c.formErrorMessage = "System error. Creating media failed"
			return nil
		}
	}

	c.formSuccessMessage = "Media saved successfully"
	return nil
}

func (c *productMediaComponent) Render() hb.TagInterface {
	htmlBytes, _ := mediaEmbed.ReadFile("media.html")
	jsBytes, _ := mediaEmbed.ReadFile("media.js")

	htmlContent := string(htmlBytes)
	jsContent := string(jsBytes)

	linksHelper := shared.NewLinksFromRequest(c.request)
	urlMediaLoad := linksHelper.ProductUpdate(map[string]string{
		"action":     "load-media",
		"product_id": c.productID,
	})

	urlMediaSave := linksHelper.ProductUpdate(map[string]string{
		"action":     "save-media",
		"product_id": c.productID,
	})

	urlMediaUpload := linksHelper.ProductUpdate(map[string]string{
		"action":     "upload-media",
		"product_id": c.productID,
	})

	initScript := `
		const productId = "` + c.productID + `";
		const urlMediaLoad = "` + urlMediaLoad + `";
		const urlMediaSave = "` + urlMediaSave + `";
		const urlMediaUpload = "` + urlMediaUpload + `";
	`

	return hb.Div().
		Child(hb.Script(initScript)).
		Child(hb.Raw(htmlContent)).
		Child(hb.Script(jsContent))
}

func (c *productMediaComponent) createMediaItem(media shopstore.MediaInterface, index int) hb.TagInterface {
	mediaType := media.GetType()
	mediaURL := media.GetURL()
	mediaTitle := media.GetTitle()

	icon := hb.I().Class("bi bi-image fs-4 text-info")
	if mediaType == shopstore.MEDIA_TYPE_VIDEO_MP4 {
		icon = hb.I().Class("bi bi-play-circle fs-4 text-primary")
	}

	title := mediaTitle
	if title == "" {
		title = "Media " + cast.ToString(index+1)
	}

	typeBadge := hb.Span().Class("badge bg-info").HTML("Image")
	if mediaType == shopstore.MEDIA_TYPE_VIDEO_MP4 {
		typeBadge = hb.Span().Class("badge bg-primary").HTML("Video")
	}

	titleDiv := hb.Div().Class("fw-bold").HTML(title)
	urlDiv := hb.Div().Class("text-muted small").HTML(mediaURL)

	mediaInfo := hb.Div().
		Class("d-flex align-items-center gap-3").
		Child(icon).
		Child(hb.Div().Child(titleDiv).Child(urlDiv))

	urlInput := hb.Input().
		Type("hidden").
		Name("media_url[" + cast.ToString(index) + "]").
		Value(mediaURL)

	typeInput := hb.Input().
		Type("hidden").
		Name("media_type[" + cast.ToString(index) + "]").
		Value(mediaType)

	deleteButton := hb.Button().
		Type("button").
		Class("btn btn-sm btn-outline-danger w-100").
		OnClick("removeMediaItem(this)").
		Child(hb.I().Class("bi bi-trash"))

	row := hb.Div().
		Class("row align-items-center g-3").
		Child(hb.Div().Class("col-md-8").Child(mediaInfo)).
		Child(hb.Div().Class("col-md-3").Child(typeBadge)).
		Child(hb.Div().Class("col-md-1").Child(deleteButton))

	cardBody := hb.Div().
		Class("card-body p-3").
		Child(row).
		Child(urlInput).
		Child(typeInput)

	card := hb.Div().
		Class("card border-0 shadow-sm").
		Child(cardBody)

	return hb.Div().
		Class("media-item mb-3").
		Child(card)
}
