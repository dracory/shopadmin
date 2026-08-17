package shopadmin

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/cachestore"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"
	_ "modernc.org/sqlite"
)

// mockRegistryImpl implements RegistryInterface for testing
type mockRegistryImpl struct {
	store      shopstore.StoreInterface
	cacheStore cachestore.StoreInterface
	logger     *slog.Logger
}

func (m *mockRegistryImpl) GetShopStore() shopstore.StoreInterface   { return m.store }
func (m *mockRegistryImpl) GetCacheStore() cachestore.StoreInterface { return m.cacheStore }
func (m *mockRegistryImpl) GetLogger() *slog.Logger                  { return m.logger }

func TestRoutes_Valid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	registry := &mockRegistryImpl{
		store:  store,
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}

	routes, err := Routes(registry, AdminOptions{
		ShopAdminURL: "/admin/shop",
	})
	if err != nil {
		t.Fatalf("Failed to create routes: %v", err)
	}
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}
}

func TestRoutes_NilRegistry(t *testing.T) {
	_, err := Routes(nil)
	if err == nil {
		t.Errorf("Expected error when registry is nil")
	}
}

func TestRoutes_Handler(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	registry := &mockRegistryImpl{
		store:  store,
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}

	routes, err := Routes(registry, AdminOptions{
		ShopAdminURL: "/admin/shop",
	})
	if err != nil {
		t.Fatalf("Failed to create routes: %v", err)
	}

	shopRoute := routes[0]
	htmlHandler := shopRoute.GetHTMLHandler()
	if htmlHandler == nil {
		t.Fatalf("Expected HTML handler on shop route")
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/shop", nil)
	rr := httptest.NewRecorder()

	html := htmlHandler(rr, req)
	body := rr.Body.String() + html
	if !strings.Contains(body, "Shop Dashboard") {
		t.Errorf("Expected body to contain 'Shop Dashboard'")
	}
}

// TestCustomerResolver_PassedThrough verifies CustomerResolver is wired
// from AdminOptions through to the order controllers
func TestCustomerResolver_PassedThrough(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	resolver := &mockCustomerResolver{
		name:  "John Doe",
		email: "john@example.com",
		ids:   []string{"cust-1"},
	}

	a, err := New(AdminOptions{
		Store:            store,
		Logger:           slog.New(slog.NewTextHandler(os.Stderr, nil)),
		CustomerResolver: resolver,
	})
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	admin := a.(*admin)
	if admin.customerResolver == nil {
		t.Errorf("Expected customerResolver to be set")
	}

	name, email := admin.customerResolver.FindByID(context.Background(), "cust-1")
	if name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", name)
	}
	if email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", email)
	}
}
