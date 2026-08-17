package shopadmin

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/shopadmin/testutils"
	_ "modernc.org/sqlite"
)

func TestNew_ValidOptions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	options := AdminOptions{
		Store:          store,
		Logger:         slog.New(slog.NewTextHandler(os.Stderr, nil)),
		AdminHomeURL:   "/admin",
		ShopAdminURL:   "/admin/shop",
		FileManagerURL: "/admin/files",
	}
	a, err := New(options)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}
	if a == nil {
		t.Errorf("Expected admin to be created, got nil")
	}
}

func TestNew_MissingStore(t *testing.T) {
	options := AdminOptions{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	a, err := New(options)
	if err == nil {
		t.Errorf("Expected error when store is missing")
	}
	if a != nil {
		t.Errorf("Expected nil admin when store is missing")
	}
	if !strings.Contains(err.Error(), ErrStoreRequired.Error()) {
		t.Errorf("Expected error to contain '%s', got '%s'", ErrStoreRequired.Error(), err.Error())
	}
}

func TestNew_MissingLogger(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	options := AdminOptions{
		Store: store,
	}
	a, err := New(options)
	if err == nil {
		t.Errorf("Expected error when logger is missing")
	}
	if a != nil {
		t.Errorf("Expected nil admin when logger is missing")
	}
	if !strings.Contains(err.Error(), ErrLoggerRequired.Error()) {
		t.Errorf("Expected error to contain '%s', got '%s'", ErrLoggerRequired.Error(), err.Error())
	}
}

func TestNew_Defaults(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	options := AdminOptions{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	a, err := New(options)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Verify defaults were set
	admin := a.(*admin)
	if admin.adminHomeURL != "/admin" {
		t.Errorf("Expected default adminHomeURL '/admin', got '%s'", admin.adminHomeURL)
	}
	if admin.shopAdminURL != "/admin/shop" {
		t.Errorf("Expected default shopAdminURL '/admin/shop', got '%s'", admin.shopAdminURL)
	}
}

func TestHandle_HomeController(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop", nil)
	rr := httptest.NewRecorder()

	a.Handle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Shop Dashboard") {
		t.Errorf("Expected body to contain 'Shop Dashboard', got: %s", body[:min(200, len(body))])
	}
}

func TestHandle_ProductManagerController(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=products", nil)
	rr := httptest.NewRecorder()

	a.Handle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Product Manager") {
		t.Errorf("Expected body to contain 'Product Manager'")
	}
}

func TestHandle_UnknownControllerFallsBackToHome(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=nonexistent", nil)
	rr := httptest.NewRecorder()

	a.Handle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Shop Dashboard") {
		t.Errorf("Expected fallback to home page with 'Shop Dashboard'")
	}
}

func TestHandle_AuthUserIDRedirect(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	a, err := New(AdminOptions{
		Store:      store,
		Logger:     slog.New(slog.NewTextHandler(os.Stderr, nil)),
		AuthUserID: func(r *http.Request) string { return "" },
	})
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop", nil)
	rr := httptest.NewRecorder()

	a.Handle(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect status 303, got %d", rr.Code)
	}
}

// mockCustomerResolver implements shared.CustomerResolverInterface for testing
type mockCustomerResolver struct {
	name  string
	email string
	ids   []string
}

func (m *mockCustomerResolver) FindByID(_ context.Context, id string) (string, string) {
	return m.name, m.email
}

func (m *mockCustomerResolver) SearchIDs(_ context.Context, name, email string) ([]string, error) {
	return m.ids, nil
}

func TestCustomerResolver_NilSafe(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Nil CustomerResolver should not panic
	a, err := New(AdminOptions{
		Store:            store,
		Logger:           slog.New(slog.NewTextHandler(os.Stderr, nil)),
		CustomerResolver: nil,
	})
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop?controller=orders", nil)
	rr := httptest.NewRecorder()

	a.Handle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 with nil CustomerResolver, got %d", rr.Code)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
