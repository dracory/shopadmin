package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/shopstore"
)

// UiBase is a base struct that implements shared.UiInterface.
// Subcontroller ui structs can embed this to get the Store(),
// CacheStore(), Logger(), CustomerResolver(), and Layout() methods
// for free, following the cmsstore/admin pattern.
type UiBase struct {
	StoreField                 shopstore.StoreInterface
	LoggerField                *slog.Logger
	CustomerResolverField      CustomerResolverInterface
	ProductViewSectionsField   []ProductSectionRenderer
	ProductUpdateSectionsField []ProductSectionRenderer
	LayoutField                func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}

func (u UiBase) Store() shopstore.StoreInterface                 { return u.StoreField }
func (u UiBase) Logger() *slog.Logger                            { return u.LoggerField }
func (u UiBase) CustomerResolver() CustomerResolverInterface     { return u.CustomerResolverField }
func (u UiBase) ProductViewSections() []ProductSectionRenderer   { return u.ProductViewSectionsField }
func (u UiBase) ProductUpdateSections() []ProductSectionRenderer { return u.ProductUpdateSectionsField }

func (u UiBase) Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	return u.LayoutField(w, r, webpageTitle, webpageHtml, options)
}

// NewUiBase creates a UiBase from a UiConfig
func NewUiBase(config UiConfig) UiBase {
	return UiBase{
		StoreField:                 config.Store,
		LoggerField:                config.Logger,
		CustomerResolverField:      config.CustomerResolver,
		ProductViewSectionsField:   config.ProductViewSections,
		ProductUpdateSectionsField: config.ProductUpdateSections,
		LayoutField:                config.Layout,
	}
}
