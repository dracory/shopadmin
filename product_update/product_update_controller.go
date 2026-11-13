package product_update

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
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
		case "load-categories":
			return u.handleLoadCategories(w, r, productID)
		case "save-categories":
			return u.handleSaveCategories(w, r, productID)
		case "load-status":
			return u.handleLoadStatus(w, r, productID)
		case "save-status":
			return u.handleSaveStatus(w, r, productID)
		case "load-variants":
			return u.handleLoadVariants(w, r, productID)
		case "save-variants":
			return u.handleSaveVariants(w, r, productID)
		case "load-notes":
			return u.handleLoadNotes(w, r, productID)
		case "save-notes":
			return u.handleSaveNotes(w, r, productID)
		case "load-details":
			return u.handleLoadDetails(w, r, productID)
		case "save-details":
			return u.handleSaveDetails(w, r, productID)
		}
	}

	if productID == "" {
		return shared.ErrorAlert("Product ID is required")
	}

	store := u.Store()
	if store == nil {
		return shared.ErrorAlert("Shop store not available")
	}

	product, err := store.ProductFindByID(r.Context(), productID)
	if err != nil {
		u.Logger().Error("Error. productUpdateController: ProductFindByID", "error", err.Error(), "product_id", productID)
		return shared.ErrorAlert("Product not found")
	}

	if product == nil {
		u.Logger().Warn("Warning. productUpdateController: ProductFindByID", "error", "Product not found", "product_id", productID)
		return shared.ErrorAlert("Product not found")
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
		_ = component.Handle(r)
		return component.Render().ToHTML()
	}

	sections := shared.ProcessProductSections(u.ProductUpdateSections(), r, product, u.Logger())

	pageContent := u.page(r, product, view, productID, sections)

	styleURLs := []string{
		"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.css",
		cdn.Notiflix_3_2_8_CSS(),
	}
	styleURLs = append(styleURLs, sections.Options.StyleURLs...)

	scriptURLs := []string{
		cdn.Jquery_3_7_1(),
		"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.js",
		cdn.Sweetalert2_10(),
		cdn.Notiflix_3_2_8(),
		cdn.VueJs_3(),
	}
	scriptURLs = append(scriptURLs, sections.Options.ScriptURLs...)

	return u.Layout(w, r, "Edit Product | Shop", pageContent.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		Styles:     sections.Options.Styles,
		StyleURLs:  styleURLs,
		Scripts:    sections.Options.Scripts,
		ScriptURLs: scriptURLs,
	})
}

func (u *ui) page(r *http.Request, product shopstore.ProductInterface, view string, productID string, sections shared.ProcessedSections) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Product Manager", URL: shared.URLR(r, shared.CONTROLLER_PRODUCTS, nil)},
		{Name: "Edit Product", URL: shared.URLR(r, shared.CONTROLLER_PRODUCT_UPDATE, map[string]string{"product_id": productID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(r, shared.CONTROLLER_PRODUCTS, nil))

	linksHelper := shared.NewLinksFromRequest(r)

	tabs := hb.UL().Class("nav nav-tabs").
		Class("mb-3").
		Child(hb.LI().Class("nav-item").
			Child(hb.Hyperlink().Class("nav-link").
				ClassIf(view == "details", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "details",
				})).
				HTML("Details"))).
		Child(hb.LI().Class("nav-item").
			Child(hb.Hyperlink().Class("nav-link").
				ClassIf(view == "media", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "media",
				})).
				HTML("Media"))).
		Child(hb.LI().Class("nav-item").
			Child(hb.Hyperlink().Class("nav-link").
				ClassIf(view == "tags", "active").
				Href(linksHelper.ProductUpdate(map[string]string{
					"product_id": productID,
					"view":       "tags",
				})).
				HTML("Tags"))).
		Child(hb.LI().Class("nav-item").
			Child(hb.Hyperlink().Class("nav-link").
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

	// All views render their own self-contained card (title + buttons
	// in the header).
	card := body

	topSection := sections.RenderPosition(shared.PositionTop)
	beforeFormSection := sections.RenderPosition(shared.PositionBeforeForm)
	afterFormSection := sections.RenderPosition(shared.PositionAfterForm)
	sidebarSection := sections.RenderPosition(shared.PositionSidebar)
	bottomSection := sections.RenderPosition(shared.PositionBottom)

	// The Status and Categories cards are managed independently through
	// their own load/save handlers. Status sits on top of the other
	// cards; Categories sits below the product details card.
	var statusCard hb.TagInterface
	var categoriesCard hb.TagInterface
	var variantsCard hb.TagInterface
	var adminNotesCard hb.TagInterface
	if view == "details" {
		statusComponent := NewProductStatusComponent(u)
		statusComponent.Mount(r, product, productID)
		statusCard = statusComponent.Render()

		categoriesComponent := NewProductCategoriesComponent(u)
		categoriesComponent.Mount(r, product, productID)
		categoriesCard = categoriesComponent.Render()

		variantsComponent := NewProductVariantsComponent(u)
		variantsComponent.Mount(r, product, productID)
		variantsCard = variantsComponent.Render()

		adminNotesComponent := NewProductAdminNotesComponent(u)
		adminNotesComponent.Mount(r, product, productID)
		adminNotesCard = adminNotesComponent.Render()
	}

	mainCardGroup := hb.Div().
		ChildIf(beforeFormSection != nil, beforeFormSection).
		ChildIf(statusCard != nil, statusCard).
		Child(card).
		ChildIf(categoriesCard != nil, categoriesCard).
		ChildIf(variantsCard != nil, variantsCard).
		ChildIf(adminNotesCard != nil, adminNotesCard).
		ChildIf(afterFormSection != nil, afterFormSection)

	var middleLayout hb.TagInterface
	if sidebarSection != nil {
		middleLayout = hb.Div().Class("row").
			Child(hb.Div().Class("col-lg-8").Child(mainCardGroup)).
			Child(hb.Div().Class("col-lg-4").Child(sidebarSection))
	} else {
		middleLayout = mainCardGroup
	}

	return hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Products. Edit Product", buttonCancel)).
		ChildIf(topSection != nil, topSection).
		Child(productTitle).
		Child(tabs).
		Child(middleLayout).
		ChildIf(bottomSection != nil, bottomSection)
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
	_, _ = w.Write(jsonBytes)
}

// writeJSONError writes a JSON error response safely (bug #7 fix).
// Previously the code did _, _ = w.Write([]byte(`{"status":"error","message":"Failed to parse upload: ` + err.Error() + `"}`))
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

	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)
	// #nosec G120 -- request body is bounded by MaxBytesReader above
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
		_ = file.Close()
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

	_ = product.SetMetas(metas)

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

	_ = product.SetMetas(metas)

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Tags saved successfully", nil)
	return ""
}

func (u *ui) handleLoadCategories(w http.ResponseWriter, r *http.Request, productID string) string {
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

	categoryQuery := shopstore.NewCategoryQuery()
	categoryQuery.SetStatus(shopstore.CATEGORY_STATUS_ACTIVE)
	categories, err := shopStore.CategoryList(ctx, categoryQuery)
	if err != nil {
		u.Logger().Error("Failed to load categories", "error", err)
		writeJSONError(w, "Failed to load categories")
		return ""
	}

	categoryItems := []map[string]any{}
	for _, category := range categories {
		categoryItems = append(categoryItems, map[string]any{
			"id":    category.GetID(),
			"title": category.GetTitle(),
		})
	}

	selected, err := shopStore.ProductCategoryIDs(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load product categories", "error", err)
		writeJSONError(w, "Failed to load product categories")
		return ""
	}

	writeJSON(w, "success", "Categories loaded successfully", map[string]any{
		"categories": categoryItems,
		"selected":   selected,
	})
	return ""
}

func (u *ui) handleSaveCategories(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		CategoryIDs []string `json:"category_ids"`
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

	err = shopStore.ProductCategorySet(ctx, productID, reqBody.CategoryIDs)
	if err != nil {
		u.Logger().Error("Failed to save product categories", "error", err)
		writeJSONError(w, "Failed to save product categories")
		return ""
	}

	writeJSON(w, "success", "Categories saved successfully", nil)
	return ""
}

// handleLoadStatus returns the product status for the Status card.
func (u *ui) handleLoadStatus(w http.ResponseWriter, r *http.Request, productID string) string {
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

	writeJSON(w, "success", "Status loaded successfully", map[string]any{
		"status_value": product.GetStatus(),
	})
	return ""
}

// handleSaveStatus updates only the product status; no other product
// data is modified.
func (u *ui) handleSaveStatus(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Status string `json:"status"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		u.Logger().Error("handleSaveStatus: failed to decode request body", "error", err, "body", string(bodyBytes))
		writeJSONError(w, "Invalid request body")
		return ""
	}

	if reqBody.Status == "" {
		writeJSONError(w, "Status is required")
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

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Status saved successfully", nil)
	return ""
}

// variantDimensionPayload is one dimension of the parent's variant
// schema (e.g. size with options S/M/L).
type variantDimensionPayload struct {
	Name     string   `json:"name"`
	Required bool     `json:"required"`
	Options  []string `json:"options"`
}

// variantPayload is one generated variant row (a dimension combination
// plus its purchasable attributes).
type variantPayload struct {
	Values   map[string]string `json:"values"`
	SKU      string            `json:"sku"`
	Price    string            `json:"price"`
	Quantity string            `json:"quantity"`
}

// variantComboKey builds a deterministic key for a variant combination,
// ordered by the dimension schema.
func variantComboKey(dimensions []variantDimensionPayload, values map[string]string) string {
	parts := make([]string, 0, len(dimensions))
	for _, dim := range dimensions {
		if v, ok := values[dim.Name]; ok {
			parts = append(parts, dim.Name+":"+v)
		}
	}
	return strings.Join(parts, "|")
}

// handleLoadVariants returns the parent's variant dimension schema and
// the existing variant products.
func (u *ui) handleLoadVariants(w http.ResponseWriter, r *http.Request, productID string) string {
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

	schemas, err := product.GetVariantMatrixSchemas()
	if err != nil {
		u.Logger().Error("Failed to parse variant schema", "error", err)
		writeJSONError(w, "Failed to load variant schema")
		return ""
	}

	dimensions := make([]variantDimensionPayload, 0, len(schemas))
	for _, s := range schemas {
		if s.Name == "" {
			continue
		}
		dimensions = append(dimensions, variantDimensionPayload{
			Name:     s.Name,
			Required: s.Required,
			Options:  s.Options,
		})
	}

	existing, err := shopStore.ProductVariantList(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load variants", "error", err)
		writeJSONError(w, "Failed to load variants")
		return ""
	}

	variants := make([]variantPayload, 0, len(existing))
	for _, v := range existing {
		values, err := v.GetVariantMatrixValues()
		if err != nil {
			u.Logger().Error("Failed to parse variant values", "error", err, "variant_id", v.GetID())
			continue
		}
		variants = append(variants, variantPayload{
			Values:   values,
			SKU:      v.GetMeta("sku"),
			Price:    v.GetPrice(),
			Quantity: v.GetQuantity(),
		})
	}

	writeJSON(w, "success", "Variants loaded successfully", map[string]any{
		"dimensions": dimensions,
		"variants":   variants,
	})
	return ""
}

// handleSaveVariants replaces the parent's variant schema and diffs the
// variant products: missing combinations are created, changed ones
// updated, removed ones soft deleted.
func (u *ui) handleSaveVariants(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Dimensions []variantDimensionPayload `json:"dimensions"`
		Variants   []variantPayload          `json:"variants"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		u.Logger().Error("handleSaveVariants: failed to decode request body", "error", err, "body", string(bodyBytes))
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

	// Update the parent's dimension schema. A nil slice marshals to
	// "null", which clears the schema entirely.
	var schemas []shopstore.VariantMatrixSchema
	for _, d := range reqBody.Dimensions {
		if d.Name == "" || len(d.Options) == 0 {
			continue
		}
		schemas = append(schemas, shopstore.VariantMatrixSchema{
			Name:     d.Name,
			Required: d.Required,
			Options:  d.Options,
		})
	}
	if err := product.SetVariantMatrixSchemas(schemas); err != nil {
		u.Logger().Error("Failed to set variant schema", "error", err)
		writeJSONError(w, "Failed to save variant schema")
		return ""
	}
	if err := shopStore.ProductUpdate(ctx, product); err != nil {
		u.Logger().Error("Failed to save product schema", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	// Index existing variants by combination key
	existing, err := shopStore.ProductVariantList(ctx, productID)
	if err != nil {
		u.Logger().Error("Failed to load variants", "error", err)
		writeJSONError(w, "Failed to load variants")
		return ""
	}
	type keyedVariant struct {
		key     string
		variant shopstore.ProductInterface
	}
	existingByKey := map[string]shopstore.ProductInterface{}
	existingKeyed := make([]keyedVariant, 0, len(existing))
	for _, v := range existing {
		values, err := v.GetVariantMatrixValues()
		if err != nil {
			continue
		}
		key := variantComboKey(reqBody.Dimensions, values)
		if key != "" {
			existingByKey[key] = v
		}
		existingKeyed = append(existingKeyed, keyedVariant{key: key, variant: v})
	}

	// Create or update every requested combination
	seen := map[string]bool{}
	for _, vp := range reqBody.Variants {
		key := variantComboKey(reqBody.Dimensions, vp.Values)
		if key == "" {
			continue
		}
		seen[key] = true

		price := vp.Price
		if price == "" {
			price = product.GetPrice() // inherit parent price
		}

		if v, ok := existingByKey[key]; ok {
			v.SetTitle(product.GetTitle())
			v.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)
			v.SetPrice(price)
			v.SetQuantity(vp.Quantity)
			_ = v.SetVariantMatrixValues(vp.Values)
			_ = v.SetMeta("sku", vp.SKU)
			if err := shopStore.ProductUpdate(ctx, v); err != nil {
				u.Logger().Error("Failed to update variant", "error", err, "variant_id", v.GetID())
				writeJSONError(w, "Failed to save variant")
				return ""
			}
			continue
		}

		variant := shopstore.NewProduct()
		variant.SetParentID(productID)
		variant.SetTitle(product.GetTitle())
		variant.SetStatus(shopstore.PRODUCT_STATUS_ACTIVE)
		variant.SetPrice(price)
		variant.SetQuantity(vp.Quantity)
		_ = variant.SetVariantMatrixValues(vp.Values)
		_ = variant.SetMeta("sku", vp.SKU)
		if err := shopStore.ProductCreate(ctx, variant); err != nil {
			u.Logger().Error("Failed to create variant", "error", err)
			writeJSONError(w, "Failed to save variant")
			return ""
		}
	}

	// Soft delete variants whose combination was removed. Variants that
	// no longer match any dimension (empty key) are always removed.
	for _, ev := range existingKeyed {
		if seen[ev.key] && ev.key != "" {
			continue
		}
		if err := shopStore.ProductSoftDeleteByID(ctx, ev.variant.GetID()); err != nil {
			u.Logger().Error("Failed to delete variant", "error", err, "variant_id", ev.variant.GetID())
		}
	}

	writeJSON(w, "success", "Variants saved successfully", nil)
	return ""
}

// handleLoadNotes returns the product memo for the Admin Notes card.
func (u *ui) handleLoadNotes(w http.ResponseWriter, r *http.Request, productID string) string {
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

	writeJSON(w, "success", "Admin notes loaded successfully", map[string]any{
		"memo": product.GetMemo(),
	})
	return ""
}

// handleSaveNotes updates only the product memo; no other product data
// is modified.
func (u *ui) handleSaveNotes(w http.ResponseWriter, r *http.Request, productID string) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		writeJSONError(w, "Shop store not available")
		return ""
	}

	var reqBody struct {
		Memo string `json:"memo"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		u.Logger().Error("handleSaveNotes: failed to decode request body", "error", err, "body", string(bodyBytes))
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

	product.SetMemo(reqBody.Memo)

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Admin notes saved successfully", nil)
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
			"title":       product.GetTitle(),
			"description": product.GetDescription(),
			"price":       product.GetPrice(),
			"quantity":    product.GetQuantity(),
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
		Title       string          `json:"title"`
		Description json.RawMessage `json:"description"`
		Price       json.RawMessage `json:"price"`
		Quantity    json.RawMessage `json:"quantity"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		u.Logger().Error("handleSaveDetails: failed to decode request body", "error", err, "body", string(bodyBytes))
		writeJSONError(w, "Invalid request body")
		return ""
	}

	// Description may arrive as a string, an object (e.g. empty {} from
	// Summernote before init), or a number. Coerce safely to string.
	description := ""
	if len(reqBody.Description) > 0 {
		var descStr string
		if err := json.Unmarshal(reqBody.Description, &descStr); err == nil {
			description = descStr
		} else {
			// Not a JSON string — use the raw bytes as a fallback string
			description = string(reqBody.Description)
		}
	}

	// Price and quantity may arrive as a JSON string or a JSON number.
	// Vue's v-model on <input type="number"> auto-casts the value to a
	// number, so coercing here prevents "Invalid request body" errors.
	price := rawMessageToString(reqBody.Price)
	quantity := rawMessageToString(reqBody.Quantity)

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

	product.SetTitle(reqBody.Title)
	product.SetDescription(description)
	product.SetPrice(price)
	product.SetQuantity(quantity)

	err = shopStore.ProductUpdate(ctx, product)
	if err != nil {
		u.Logger().Error("Failed to save product", "error", err)
		writeJSONError(w, "Failed to save product")
		return ""
	}

	writeJSON(w, "success", "Details saved successfully", nil)
	return ""
}

// rawMessageToString coerces a JSON value to a string. It accepts both
// JSON strings ("25") and JSON numbers (25), returning "" for empty input.
func rawMessageToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return string(raw)
}
