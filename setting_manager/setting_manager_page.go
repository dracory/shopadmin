package setting_manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/shopadmin/shared"
)

var (
	//go:embed settings.html
	settingsHTML string

	//go:embed settings.js
	settingsJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Shop store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{Name: "Home", URL: shared.AdminHomeURL(r)},
		{Name: "Shop", URL: shared.URLR(r, shared.CONTROLLER_HOME, nil)},
		{Name: "Settings", URL: shared.URLR(r, shared.CONTROLLER_SETTINGS, nil)},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadSettings": linksHelper.Settings(map[string]string{"action": actionLoadSettings}),
		"saveSetting":  linksHelper.Settings(map[string]string{"action": actionSettingSave}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.settingManagerUrls = %s;", string(urlsJSON)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Shop. Settings")).
		Child(urlsScript).
		Child(hb.Raw(settingsHTML)).
		Child(hb.Script(settingsJS))

	return u.Layout(w, r, "Settings | Shop", content.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Sweetalert2_10(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
