package discount_update

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

const (
	actionLoadDetails = "load-details"
	actionSaveDetails = "save-details"
)

// UiInterface defines the discount update controller's UI interface
type UiInterface interface {
	shared.UiInterface
	DiscountUpdate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new discount update controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// DiscountUpdate handles the discount update controller requests
func (u *ui) DiscountUpdate(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the discount update request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")
	discountID := req.GetStringTrimmed(r, "discount_id")

	if action != "" {
		switch action {
		case actionLoadDetails:
			return u.handleLoadDetails(w, r, discountID)
		case actionSaveDetails:
			return u.handleSaveDetails(w, r, discountID)
		}
	}

	if discountID == "" {
		return shared.ErrorAlert("Discount ID is required")
	}

	store := u.Store()
	if store == nil {
		return shared.ErrorAlert("Shop store not available")
	}

	discount, err := store.DiscountFindByID(r.Context(), discountID)
	if err != nil {
		u.Logger().Error("discountUpdateController: DiscountFindByID", "error", err.Error(), "discount_id", discountID)
		return shared.ErrorAlert("Discount not found")
	}

	if discount == nil {
		u.Logger().Warn("discountUpdateController: DiscountFindByID", "error", "Discount not found", "discount_id", discountID)
		return shared.ErrorAlert("Discount not found")
	}

	pageContent := u.page(r, discount, discountID)

	return u.Layout(w, r, "Edit Discount | Shop", pageContent.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.VueElementPlusJs_2_13_7(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.VueElementPlusCss_2_13_7(),
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}

func (u *ui) page(r *http.Request, discount shopstore.DiscountInterface, discountID string) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Discounts", URL: shared.URLR(r, shared.CONTROLLER_DISCOUNTS, nil)},
		{Name: "Edit Discount", URL: shared.URLR(r, shared.CONTROLLER_DISCOUNT_UPDATE, map[string]string{"discount_id": discountID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(r, shared.CONTROLLER_DISCOUNTS, nil))

	heading := hb.Heading1().
		HTML("Shop. Discount. Edit Discount").
		Child(buttonCancel)

	discountTitle := hb.Heading2().
		Class("mb-3").
		Text("Discount: ").
		Text(discount.GetTitle())

	component := NewDiscountDetailsComponent(u)
	component.Mount(r, discount, discountID)
	body := component.Render()

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Child(hb.Heading4().HTML("Discount Details").Style("margin-bottom:0;display:inline-block;")),
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
		Child(discountTitle).
		Child(card)
}

func (u *ui) handleLoadDetails(w http.ResponseWriter, r *http.Request, discountID string) string {
	w.Header().Set("Content-Type", "application/json")
	shopStore := u.Store()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	discount, err := shopStore.DiscountFindByID(r.Context(), discountID)
	if err != nil {
		u.Logger().Error("Failed to load discount", "error", err)
		return api.Error("Failed to load discount").ToString()
	}

	if discount == nil {
		return api.Error("Discount not found").ToString()
	}

	return api.SuccessWithData("Details loaded successfully", map[string]any{
		"code":                  discount.GetCode(),
		"title":                 discount.GetTitle(),
		"description":           discount.GetDescription(),
		"type":                  discount.GetType(),
		"amount":                cast.ToFloat64(discount.GetAmount()),
		"status":                discount.GetStatus(),
		"starts_at":             discount.GetStartsAtCarbon().ToDateTimeString(carbon.UTC),
		"ends_at":               discount.GetEndsAtCarbon().ToDateTimeString(carbon.UTC),
		"memo":                  discount.GetMemo(),
		"max_uses":              discount.GetMaxUses(),
		"max_uses_count":        discount.GetMaxUsesCount(),
		"max_uses_per_customer": discount.GetMaxUsesPerCustomer(),
	}).ToString()
}

func (u *ui) handleSaveDetails(w http.ResponseWriter, r *http.Request, discountID string) string {
	w.Header().Set("Content-Type", "application/json")
	shopStore := u.Store()
	if shopStore == nil {
		return api.Error("Shop store not available").ToString()
	}

	var reqBody struct {
		Code               string  `json:"code"`
		Title              string  `json:"title"`
		Description        string  `json:"description"`
		Type               string  `json:"type"`
		Amount             float64 `json:"amount"`
		Status             string  `json:"status"`
		StartsAt           string  `json:"starts_at"`
		EndsAt             string  `json:"ends_at"`
		Memo               string  `json:"memo"`
		MaxUses            int     `json:"max_uses"`
		MaxUsesPerCustomer int     `json:"max_uses_per_customer"`
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	u.Logger().Debug("discount_update handleSaveDetails raw body", "body", string(bodyBytes), "discount_id", discountID)
	json.Unmarshal(bodyBytes, &reqBody)

	discount, err := shopStore.DiscountFindByID(r.Context(), discountID)
	if err != nil {
		u.Logger().Error("Failed to load discount", "error", err)
		return api.Error("Failed to load discount").ToString()
	}

	if discount == nil {
		return api.Error("Discount not found").ToString()
	}

	// Enforce code uniqueness: if the code is changing, ensure no other discount
	// already uses it. DiscountFindByCode is the checkout lookup path, so a
	// duplicate would silently return one of two codes.
	trimmedCode := strings.TrimSpace(reqBody.Code)
	if trimmedCode != "" && trimmedCode != discount.GetCode() {
		existing, err := shopStore.DiscountFindByCode(r.Context(), trimmedCode)
		if err == nil && existing != nil && existing.GetID() != discountID {
			return api.Error("Discount code already in use by another discount").ToString()
		}
	}

	startsAt := reqBody.StartsAt
	if strings.TrimSpace(startsAt) == "" {
		startsAt = carbon.Now().SubDays(1).ToDateTimeString(carbon.UTC)
	}

	endsAt := reqBody.EndsAt
	if strings.TrimSpace(endsAt) == "" {
		endsAt = carbon.Now().AddYears(10).ToDateTimeString(carbon.UTC)
	}

	discount.SetCode(reqBody.Code).
		SetTitle(reqBody.Title).
		SetDescription(reqBody.Description).
		SetType(reqBody.Type).
		SetAmount(reqBody.Amount).
		SetStatus(reqBody.Status).
		SetStartsAt(carbon.Parse(startsAt).ToDateTimeString(carbon.UTC)).
		SetEndsAt(carbon.Parse(endsAt).ToDateTimeString(carbon.UTC)).
		SetMemo(reqBody.Memo).
		SetMaxUses(reqBody.MaxUses).
		SetMaxUsesPerCustomer(reqBody.MaxUsesPerCustomer)

	err = shopStore.DiscountUpdate(r.Context(), discount)
	if err != nil {
		u.Logger().Error("Failed to save discount", "error", err.Error(), "discount_id", discountID, "type", reqBody.Type, "status", reqBody.Status)
		return api.ErrorWithData("Failed to save discount", map[string]any{"discount_id": discountID}).ToString()
	}

	return api.Success("Details saved successfully").ToString()
}
