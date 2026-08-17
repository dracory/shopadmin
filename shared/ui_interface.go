package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/shopstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
// This follows the cmsstore/admin pattern.
//
// Customer resolution is via CustomerResolverInterface rather than a
// userstore dependency.
type UiInterface interface {
	Store() shopstore.StoreInterface
	Logger() *slog.Logger

	// CustomerResolver returns the CustomerResolverInterface, or nil.
	CustomerResolver() CustomerResolverInterface

	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}
