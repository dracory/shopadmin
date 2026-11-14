package setting_manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSettingManagerPageRenders(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=settings", http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.SettingManager(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Settings") {
		t.Fatal("expected page heading 'Settings'")
	}
}
