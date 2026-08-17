// Package shopadmin provides a standalone shop admin interface following
// the folder-per-controller pattern. Each controller is in its own
// subfolder and handles its own views and AJAX data.
//
// This module is modeled on github.com/dracory/cmsstore/admin.
package shopadmin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/cachestore"
	"github.com/dracory/req"
	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopstore"
)

// AdminOptions contains all dependencies and configuration for the shop admin.
//
// Store, CacheStore, and Logger replace the in-repo version's
// Registry field (which was app.AppInterface). This matches the
// cmsstore/admin convention where stores are passed directly.
//
// Customer resolution is via CustomerResolverInterface rather than a
// userstore dependency, keeping shopadmin decoupled from any specific
// auth/user-management package.
type AdminOptions struct {
	// Store is the shopstore.StoreInterface (required)
	Store shopstore.StoreInterface

	// CacheStore is optional — used for flash messages.
	// If nil, flash errors fall back to inline HTML rendering.
	CacheStore cachestore.StoreInterface

	// Logger is required (matches cmsstore requirement)
	Logger *slog.Logger

	// CustomerResolver resolves customer data for order views.
	// Optional — nil means customer fields stay empty and customer
	// filtering is disabled. Called by order_details and order_manager.
	CustomerResolver CustomerResolverInterface

	// FuncLayout is an optional function to render the admin interface
	// inside your own layout (branding, menus, etc.).
	// If nil, a default bare-bones HTML page is used (Bootstrap + Vue CDN).
	// Uses anonymous struct to match cmsstore/admin exactly, so consumers
	// can reuse their cmsstore layout function for shopadmin.
	FuncLayout func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string

	// AdminHomeURL is the URL for the admin home page (default: "/admin")
	AdminHomeURL string

	// ShopAdminURL is the base URL for shop admin (default: "/admin/shop")
	ShopAdminURL string

	// FileManagerURL is the URL for the file manager (optional)
	FileManagerURL string

	// AuthUserID returns the authenticated user ID from the request.
	// If it returns "", the user is treated as unauthenticated.
	AuthUserID func(r *http.Request) string
}

// AdminInterface defines the interface for the shop admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	store            shopstore.StoreInterface
	cacheStore       cachestore.StoreInterface
	logger           *slog.Logger
	customerResolver CustomerResolverInterface
	funcLayout       func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	adminHomeURL   string
	shopAdminURL   string
	fileManagerURL string
	authUserID     func(r *http.Request) string
	routes         map[string]func(w http.ResponseWriter, r *http.Request)
}

// New creates a new shop admin instance.
// Returns ErrStoreRequired if Store is nil, ErrLoggerRequired if Logger is nil.
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Store == nil {
		return nil, ErrStoreRequired
	}
	if opts.Logger == nil {
		return nil, ErrLoggerRequired
	}

	// Set defaults (fixes pre-existing bug #12 where AdminHomeURL had no default)
	if opts.AdminHomeURL == "" {
		opts.AdminHomeURL = "/admin"
	}
	if opts.ShopAdminURL == "" {
		opts.ShopAdminURL = "/admin/shop"
	}

	a := &admin{
		store:            opts.Store,
		cacheStore:       opts.CacheStore,
		logger:           opts.Logger,
		customerResolver: opts.CustomerResolver,
		funcLayout:       opts.FuncLayout,
		adminHomeURL:     opts.AdminHomeURL,
		shopAdminURL:     opts.ShopAdminURL,
		fileManagerURL:   opts.FileManagerURL,
		authUserID:       opts.AuthUserID,
	}

	// Build routes once at construction time (fixes perf issue #4
	// where the in-repo version called Routes() on every request)
	a.routes = a.buildRoutes()

	return a, nil
}

// Handle processes all shop admin requests.
// Config values are injected into the request context (following the
// cmsstore/admin pattern). Route lookup is map-based (fixes bug #13
// where HasPrefix matching was ambiguous).
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if a.authUserID != nil && a.authUserID(r) == "" {
		http.Redirect(w, r, a.adminHomeURL, http.StatusSeeOther)
		return
	}

	// Inject config into request context (like cmsstore/admin)
	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.adminHomeURL)
	ctx = context.WithValue(ctx, shared.KeyShopAdminURL, a.shopAdminURL)
	ctx = context.WithValue(ctx, shared.KeyFileManagerURL, a.fileManagerURL)
	r = r.WithContext(ctx)

	// Map-based route lookup (not HasPrefix — fixes bug #13)
	controller := req.GetStringTrimmed(r, "controller")
	if controller == "" {
		controller = shared.CONTROLLER_HOME
	}

	handler, ok := a.routes[controller]
	if !ok {
		handler = a.routes[shared.CONTROLLER_HOME]
	}

	handler(w, r)
}

// buildRoutes creates the controller dispatch map once at construction time.
func (a *admin) buildRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	uiConfig := shared.UiConfig{
		Store:            a.store,
		CacheStore:       a.cacheStore,
		Logger:           a.logger,
		CustomerResolver: a.customerResolver,
		Layout:           a.render,
	}

	return buildControllerRoutes(uiConfig, a.fileManagerURL)
}

// render wraps content in the layout. If FuncLayout is provided and
// returns non-empty, it is used; otherwise the default shared.Layout
// is used (following the cmsstore/admin pattern).
//
// When FuncLayout is set, the default shared.Layout is NOT computed
// (avoids wasted work — fixes #14).
func (a *admin) render(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	// If a custom layout is provided, try it first (fixes #14 — avoids
	// computing the default layout when FuncLayout is set)
	if a.funcLayout != nil {
		custom := a.funcLayout(webpageTitle, webpageHtml, options)
		if custom != "" {
			if w != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(custom))
				return ""
			}
			return custom
		}
	}

	webpage := shared.Layout(w, r, webpageTitle, webpageHtml, options)

	// w may be nil when a subcontroller calls Layout() to get the HTML
	// string without writing to the response (e.g. home controller).
	if w != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(webpage))
		return ""
	}

	return webpage
}
