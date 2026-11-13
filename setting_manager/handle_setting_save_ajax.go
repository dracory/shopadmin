package setting_manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/api"
)

func (u *ui) handleSettingSave(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	shopStore := u.Store()
	if shopStore == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Shop store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	reqBody.Key = strings.TrimSpace(reqBody.Key)
	if reqBody.Key == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Setting key is required").ToString()))
		return ""
	}

	if err := shopStore.SettingSet(ctx, reqBody.Key, reqBody.Value); err != nil {
		u.Logger().Error("Failed to save setting", "error", err, "key", reqBody.Key)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(api.Error("Failed to save setting").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(api.Success("Setting saved successfully").ToString()))
	return ""
}
