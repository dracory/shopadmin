package setting_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleSettingDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		BulkSettingIDs []string `json:"bulk_setting_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.BulkSettingIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("No setting IDs provided").ToString()))
		return ""
	}

	for _, settingID := range reqBody.BulkSettingIDs {
		if err := shopStore.SettingDeleteByID(ctx, settingID); err != nil {
			u.Logger().Error("Failed to delete setting", "error", err, "setting_id", settingID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Settings deleted successfully").ToString()))
	return ""
}
