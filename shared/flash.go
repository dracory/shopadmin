package shared

import (
	"fmt"
	"net/http"

	"github.com/dracory/cachestore"
	"github.com/dracory/hb"
	"github.com/dracory/uid"
)

// ToFlashError stores an error message in the cache store and redirects
// to a flash page URL. If the cache store is nil, it renders the error
// inline as HTML instead. The writer is nil-checked (fixes pre-existing
// bug #14 where controllers passed nil as the writer).
//
// Parameters:
//   - cacheStore: optional cache store for flash message storage
//   - w: optional response writer for redirect (may be nil)
//   - r: the HTTP request
//   - message: the error message to display
//   - redirectURL: where to redirect after showing the flash
//   - seconds: how many seconds to show the flash (unused if cacheStore is nil)
//
// Returns the HTML to render if no redirect is possible, or empty string
// if a redirect was issued.
func ToFlashError(cacheStore cachestore.StoreInterface, w http.ResponseWriter, r *http.Request, message, redirectURL string, seconds int) string {
	if cacheStore != nil && w != nil {
		messageID := uid.HumanUid()
		// Use the seconds parameter as the TTL (fixes #10 — previously
		// hardcoded to 0, which either never expires or expires immediately
		// depending on cachestore semantics)
		cacheStore.Set(messageID, message, int64(seconds))
		flashURL := fmt.Sprintf("%s/flash?message_id=%s", AdminHomeURL(r), messageID)
		http.Redirect(w, r, flashURL, http.StatusSeeOther)
		return ""
	}

	// Fallback: render the error inline as HTML
	return hb.Div().
		Class("alert alert-danger").
		Style("margin: 20px;").
		HTML(message).
		ToHTML()
}
