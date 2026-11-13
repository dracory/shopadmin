package setting_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleSettingDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		SettingID string `json:"setting_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.SettingID == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Setting ID is required").ToString()))
		return ""
	}

	if err := shopStore.SettingDeleteByID(ctx, reqBody.SettingID); err != nil {
		u.Logger().Error("Failed to delete setting", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to delete setting").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Setting deleted successfully").ToString()))
	return ""
}
