package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// uiInterfaceImpl is a minimal UiInterface implementation used to verify
// the interface contract is implementable.
type uiInterfaceImpl struct {
	UiBase
}

func TestUiInterfaceContract(t *testing.T) {
	var iface UiInterface = &uiInterfaceImpl{UiBase: NewUiBase(UiConfig{
		Layout: func(w http.ResponseWriter, r *http.Request, title, html string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return title + html
		},
	})}

	if iface.Store() != nil || iface.Logger() != nil || iface.CustomerResolver() != nil {
		t.Fatal("expected nil dependencies")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	if got := iface.Layout(w, r, "a", "b", struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{}); got != "ab" {
		t.Fatalf("unexpected layout output %q", got)
	}
}
