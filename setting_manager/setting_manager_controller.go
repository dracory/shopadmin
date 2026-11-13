package setting_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
)

const (
	actionLoadSettings          = "load-settings"
	actionSettingSave           = "save-setting"
	actionSettingDelete         = "delete-setting"
	actionSettingDeleteSelected = "delete-selected"
)

// UiInterface defines the setting manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	SettingManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new setting manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// SettingManager handles the setting manager controller requests
func (u *ui) SettingManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the setting manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadSettings:
		return u.handleLoadSettings(w, r)
	case actionSettingSave:
		return u.handleSettingSave(w, r)
	case actionSettingDelete:
		return u.handleSettingDelete(w, r)
	case actionSettingDeleteSelected:
		return u.handleSettingDeleteSelected(w, r)
	default:
		return u.renderPage(w, r)
	}
}
