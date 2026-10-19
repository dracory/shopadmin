package shared

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/shopstore"
)

// CustomerResolverInterface resolves customer data for order views.
// The host project provides an implementation — shopadmin does not
// care where the data comes from (userstore, CRM, external API, etc.).
//
// All methods are optional — a nil implementation means customer
// fields stay empty and customer filtering is disabled.
//
// This interface replaces the former userstore dependency, keeping
// shopadmin decoupled from any specific auth/user-management package.
type CustomerResolverInterface interface {
	// FindByID returns customer display name and email for a given ID.
	// Returns empty strings if not found — no error for "not found".
	// Called by order_details and order_manager controllers.
	FindByID(ctx context.Context, customerID string) (name, email string)

	// SearchIDs returns customer IDs matching the given name and/or
	// email substrings. Empty string means "no filter on that field".
	// Called by order_manager controller.
	SearchIDs(ctx context.Context, name, email string) ([]string, error)
}

// UiConfig holds the dependencies passed to subcontroller UI factories.
// The Layout function uses an anonymous struct for options to match
// cmsstore/admin exactly, allowing consumers to reuse their cmsstore
// layout function for shopadmin.
//
// Customer resolution is via CustomerResolverInterface rather than a
// userstore dependency, keeping shopadmin decoupled from any specific
// auth/user-management package.
// ProductSectionPosition represents placement location for custom sections.
type ProductSectionPosition string

const (
	PositionTop        ProductSectionPosition = "top"
	PositionBeforeForm ProductSectionPosition = "before_form"
	PositionAfterForm  ProductSectionPosition = "after_form"
	PositionSidebar    ProductSectionPosition = "sidebar"
	PositionBottom     ProductSectionPosition = "bottom"
)

// ProductSection represents custom markup, CSS, JS, and external assets.
type ProductSection struct {
	ID       string
	Position ProductSectionPosition
	HTML     string
	CSS      string
	JS       string
	CSSURLs  []string
	JSURLs   []string
}

// ProductSectionRenderer is a callback that generates a ProductSection.
type ProductSectionRenderer func(
	r *http.Request,
	product shopstore.ProductInterface,
) (ProductSection, error)

type UiConfig struct {
	Store                 shopstore.StoreInterface
	Logger                *slog.Logger
	CustomerResolver      CustomerResolverInterface
	ProductViewSections   []ProductSectionRenderer
	ProductUpdateSections []ProductSectionRenderer
	Layout                func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
