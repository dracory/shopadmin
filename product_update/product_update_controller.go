package product_update

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

// UiInterface defines the product update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	ProductUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
	fileManagerURL string
}

// UI creates a new product update controller UI from the given config
func UI(config shared.UiConfig, fileManagerURL string) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config), fileManagerURL: fileManagerURL}
}

// ProductUpdate handles the product update controller requests
func (u *ui) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the product update request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")
	productID := req.GetStringTrimmed(r, "product_id")
	view := req.GetStringTrimmedOr(r, "view", "details")

	// Handle AJAX actions
	if action != "" {
		switch action {
		case "load-media":
			return u.handleLoadMedia(w, r, productID)
		case "save-media":
			return u.handleSaveMedia(w, r, productID)
		case "upload-media":
			return u.handleUploadMedia(w, r, productID)
		case "load-metadata":
			return u.handleLoadMetadata(w, r, productID)
		case "save-metadata":
			return u.handleSaveMetadata(w, r, productID)
		case "load-tags":
			return u.handleLoadTags(w, r, productID)
		case "save-tags":
			return u.handleSaveTags(w, r, productID)
		case "load-details":
			return u.handleLoadDetails(w, r, productID)
		case "save-details":
			return u.handleSaveDetails(w, r, productID)
		}
	}

	if productID == "" {
		return shared.ToFlashError(u.CacheStore(), w, r, "Product ID is required", shared.AdminHomeURL(r), 10)
	}

	store := u.Store()
	if store == nil {
		return shared.ToFlashError(u.CacheStore(), w, r, "Shop store not available", shared.AdminHomeURL(r), 10)
	}

	product, err := store.ProductFindByID(r.Context(), productID)
	if err != nil {
		u.Logger().Error("Error. productUpdateController: ProductFindByID", "error", err.Error(), "product_id", productID)
		return shared.ToFlashError(u.CacheStore(), w, r, "Product not found", shared.AdminHomeURL(r), 10)
	}

	if product == nil {
		u.Logger().Warn("Warning. productUpdateController: ProductFindByID", "error", "Product not found", "product_id", productID)
		return shared.ToFlashError(u.CacheStore(), w, r, "Product not found", shared.AdminHomeURL(r), 10)
	}

	// Handle POST requests for each view
	if r.Method == http.MethodPost {
		var component interface {
			Mount(*http.Request, shopstore.ProductInterface, string)
			Handle(*http.Request) error
			Render() hb.TagInterface
		}

		switch view {
		case "details":
			component = NewProductDetailsComponent(u)
		case "metadata":
			component = NewProductMetadataComponent(u)
		case "media":
			component = NewProductMediaComponent(u)
		case "tags":
			component = NewProductTagsComponent(u)
		default:
			component = NewProductDetailsComponent(u)
		}

		component.Mount(r, product, productID)
		component.Handle(r)
		return component.Render().ToHTML()
	}

	pageContent := u.page(r, product, view, productID)

	return u.Layout(w, r, "Edit Product | Shop", pageContent.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Jquery_3_7_1(),
			"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.js",
			cdn.Sweetalert2_10(),
			"https://unpkg.com/vue@3/dist/vue.global.js",
		},
		StyleURLs: []string{
			"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.css",
		},
	})
}

func (u *ui) page(r *http.Request, product shopstore.ProductInterface, view string, productID string) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Product Manager", URL: shared.URLR(r, shared.CONTROLLER_PRODUCTS, nil)},
		{Name: "Edit Product", URL: shared.URLR(r, shared.CONTROLLER_PRODUCT_UPDATE, map[string]string{"product_id": productID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(r, shared.CONTROLLER_PRODUCTS, nil))

	heading := hb.Heading1().
		HTML("Shop. Product. Edit Product").
		Child(buttonCancel)

	linksHelper := shared.NewLinksFromRequest(r)

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "details", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "details",
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "media", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "media",
				})).
				HTML("Media"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "tags", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "tags",
				})).
				HTML("Tags"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "metadata", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "metadata",
				})).
				HTML("Metadata")))

	productTitle := hb.Heading2().
		Class("mb-3").
		Text("Product: ").
		Text(product.GetTitle())

	var body hb.TagInterface

	switch view {
	case "details":
		component := NewProductDetailsComponent(u)
		component.Mount(r, product, productID)
		body = component.Render()
	case "media":
		component := NewProductMediaComponent(u)
		component.Mount(r, product, productID)
		body = component.Render()
	case "tags":
		component := NewProductTagsComponent(u)
		component.Mount(r, product, productID)
		body = component.Render()
	case "metadata":
		component := NewProductMetadataComponent(u)
		component.Mount(r, product, productID)
		body = component.Render()
	default:
		component := NewProductDetailsComponent(u)
		component.Mount(r, product, productID)
		body = component.Render()
	}

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Child(hb.Heading4().
					HTMLIf(view == "details", "Product Details").
					HTMLIf(view == "media", "Product Media").
					HTMLIf(view == "tags", "Product Tags").
					HTMLIf(view == "metadata", "Product Metadata").
					Style("margin-bottom:0;display:inline-block;")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(body),
		)

	return hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(productTitle).
		Child(tabs).
		Child(card)
}

// writeJSON safely writes a JSON response, fixing bug #7 where error
// messages were injected into raw JSON strings (which breaks on quotes).
func writeJSON(w http.ResponseWriter, status, message string, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"status":  status,
		"message": message,
	}
	for k, v := range data {
		resp[k] = v
	}
	jsonBytes, _ := json.Marshal(resp)
	w.Write(jsonBytes)
}

// writeJSONError writes a JSON error response safely (bug #7 fix).
// Previously the code did w.Write([]byte(`{"status":"error","message":"Failed to parse upload: ` + err.Error() + `"}`))
// which breaks if err.Error() contains a double quote.
func writeJSONError(w http.ResponseWriter, message string) {
	writeJSON(w, "error", message, nil)
}

func (u *ui) handleLoadMedia(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	medias, err := shopStore.MediaList(ctx, mediaQuery)
	if err != nil {
		u.Logger().Error("Failed to load media", "error", err)
		writeJSONError(w, "Failed to load media")
		return ""
	}

	mediaItems := []map[string]any{}
	for i, media := range medias {
		mediaItems = append(mediaItems, map[string]any{
			"id":       media.GetID(),
			"fileName": media.GetTitle(),
			"url":      media.GetURL(),
			"isMain":   i == 0,
			"type":     media.GetType(),
		})
	}

	writeJSON(w, "success", "Media loaded successfully", map[string]any{
		"data": map[string]any{
			"media": mediaItems,
		},
	})
	return ""
}

// handleSaveMedia fixes bug #6: unchecked type assertions on map[string]any
// values (item["url"].(string), item["fileName"].(string), item["type"].(string))
// which would panic if the keys were missing or had non-string values.
func (u *ui) handleSaveMedia(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Media []map[string]any `json:"media"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		writeJSONError(w, "Invalid request body")
		return ""
	}

	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	existingMedias, err := shopStore.MediaList(ctx, mediaQuery)
	if err != nil {
		u.Logger().Error("Failed to load existing media", "error", err)
		writeJSONError(w, "Failed to load existing media")
		return ""
	}

	for _, existingMedia := range existingMedias {
		err := shopStore.MediaDelete(ctx, existingMedia)
		if err != nil {
			u.Logger().Error("Failed to delete media", "error", err)
		}
	}

	for i, item := range reqBody.Media {
		// Bug #6 fix: use safe type assertions instead of unchecked .(string)
		url, _ := item["url"].(string)
		fileName, _ := item["fileName"].(string)
		mediaType, _ := item["type"].(string)

		if url == "" {
			continue
		}

		media := shopstore.NewMedia()
		media.SetEntityID(productID)
		media.SetURL(url)
		media.SetTitle(fileName)
		media.SetType(mediaType)
		media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
		media.SetSequence(i)

		err := shopStore.MediaCreate(ctx, media)
		if err != nil {
			u.Logger().Error("Failed to create media", "error", err)
		}
	}

	writeJSON(w, "success", "Media saved successfully", nil)
	return ""
}

func (u *ui) handleUploadMedia(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		// Bug #7 fix: use writeJSONError instead of injecting err.Error() into raw JSON
		writeJSONError(w, fmt.Sprintf("Failed to parse upload: %s", err.Error()))
		return ""
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		files = r.MultipartForm.File["upload_file"]
	}
	if len(files) == 0 {
		writeJSONError(w, "No files uploaded")
		return ""
	}

	mediaQuery := shopstore.NewMediaQuery()
	mediaQuery.SetEntityID(productID)
	mediaQuery.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	existingMedias, _ := shopStore.MediaList(ctx, mediaQuery)
	startSequence := len(existingMedias)

	uploaded := []map[string]any{}

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			writeJSONError(w, fmt.Sprintf("Failed to open file: %s", err.Error()))
			return ""
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			writeJSONError(w, fmt.Sprintf("Failed to read file: %s", err.Error()))
			return ""
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		dataURI := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data)

		media := shopstore.NewMedia()
		media.SetEntityID(productID)
		media.SetTitle(fileHeader.Filename)
		media.SetURL(dataURI)
		media.SetType(contentType)
		media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
		media.SetSequence(startSequence + i)

		err = shopStore.MediaCreate(ctx, media)
		if err != nil {
			u.Logger().Error("Failed to create media", "error", err)
			writeJSONError(w, fmt.Sprintf("Failed to save file record: %s", err.Error()))
			return ""
		}

		uploaded = append(uploaded, map[string]any{
			"id":       media.GetID(),
			"fileName": media.GetTitle(),
			"url":      media.GetURL(),
			"type":     media.GetType(),
			"sequence": media.GetSequence(),
		})
	}

	writeJSON(w, "success", "Files uploaded successfully", map[string]any{
		"data": map[string]any{
			"media": uploaded,
		},
	})
	return ""
}

func (u *ui) handleLoadMetadata(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		writeJSONError(w, "Failed to load product")
		return ""
	}

	if product == nil {
		writeJSONError(w, "Product not found")
		return ""
	}

	metas, _ := product.GetMetas()

	metadataItems := []map[string]any{}
	for key, value := range metas {
		metadataItems = append(metadataItems, map[string]any{
			"id":    key,
			"key":   key,
			"value": value,
		})
	}

	writeJSON(w, "success", "Metadata loaded successfully", map[string]any{
		"metadata": metadataItems,
	})
	return ""
}

func (u *ui) handleSaveMetadata(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Metadata []map[string]any `json:"metadata"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		writeJSONError(w, "Invalid request body")
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		writeJSONError(w, "Failed to load product")
		return ""
	}

	if product == nil {
		writeJSONError(w, "Product not found")
		return ""
	}

	metas := make(map[string]string)
	for _, item := range reqBody.Metadata {
		if key, ok := item["key"].(string); ok {
			if value, ok := item["value"].(string); ok {
				metas[key] = value
			}
		}
	}

	product.SetMetas(metas)

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Metadata saved successfully", nil)
	return ""
}

func (u *ui) handleLoadTags(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		writeJSONError(w, "Failed to load product")
		return ""
	}

	if product == nil {
		writeJSONError(w, "Product not found")
		return ""
	}

	metas, _ := product.GetMetas()

	var tags []string
	if tagsMeta, exists := metas["tags"]; exists && tagsMeta != "" {
		tagParts := strings.Split(tagsMeta, ",")
		for _, tag := range tagParts {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	writeJSON(w, "success", "Tags loaded successfully", map[string]any{
		"tags": tags,
	})
	return ""
}

func (u *ui) handleSaveTags(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Tags []string `json:"tags"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		writeJSONError(w, "Invalid request body")
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		writeJSONError(w, "Failed to load product")
		return ""
	}

	if product == nil {
		writeJSONError(w, "Product not found")
		return ""
	}

	metas, _ := product.GetMetas()
	if metas == nil {
		metas = make(map[string]string)
	}

	tagsString := strings.Join(reqBody.Tags, ",")

	if tagsString != "" {
		metas["tags"] = tagsString
	} else {
		delete(metas, "tags")
	}

	product.SetMetas(metas)

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Tags saved successfully", nil)
	return ""
}

func (u *ui) handleLoadDetails(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		writeJSONError(w, "Failed to load product")
		return ""
	}

	if product == nil {
		writeJSONError(w, "Product not found")
		return ""
	}

	writeJSON(w, "success", "Details loaded successfully", map[string]any{
		"data": map[string]any{
			"status":      product.GetStatus(),
			"title":       product.GetTitle(),
			"description": product.GetDescription(),
			"price":       product.GetPrice(),
			"quantity":    product.GetQuantity(),
			"memo":        product.GetMemo(),
		},
	})
	return ""
}

func (u *ui) handleSaveDetails(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Status      string `json:"status"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Quantity    string `json:"quantity"`
		Memo        string `json:"memo"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		writeJSONError(w, "Invalid request body")
		return ""
	}

	product, err := shopStore.ProductFindByID(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product", "error", err)
		writeJSONError(w, "Failed to load product")
		return ""
	}

	if product == nil {
		writeJSONError(w, "Product not found")
		return ""
	}

	product.SetStatus(reqBody.Status)
	product.SetTitle(reqBody.Title)
	product.SetDescription(reqBody.Description)
	product.SetPrice(reqBody.Price)
	product.SetQuantity(reqBody.Quantity)
	product.SetMemo(reqBody.Memo)

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Details saved successfully", nil)
	return ""
}
