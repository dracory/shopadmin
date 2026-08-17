package shopadmin

import (
	"errors"
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/rtr"
	"github.com/dracory/shopadmin/shared"
)

// Routes returns the routes for the shop admin, for integration with
// the host project's router. The signature preserves the original
// Routes(registry, opts...) form — registry is RegistryInterface
// (structurally compatible with project/internal/app.AppInterface),
// so existing call sites like shopadmin.Routes(app, opts) work unchanged.
func Routes(registry RegistryInterface, opts ...AdminOptions) ([]rtr.RouteInterface, error) {
	var options AdminOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	if registry == nil {
		return nil, errors.New("registry cannot be nil")
	}

	// Build AdminOptions from registry + opts
	fullOpts := AdminOptions{
		Store:            registry.GetShopStore(),
		Logger:           registry.GetLogger(),
		CustomerResolver: options.CustomerResolver,
		FuncLayout:       options.FuncLayout,
		AdminHomeURL:     options.AdminHomeURL,
		ShopAdminURL:     options.ShopAdminURL,
		FileManagerURL:   options.FileManagerURL,
		AuthUserID:       options.AuthUserID,
	}

	// Validate store (fixes #11 — New() validates but Routes() didn't)
	if fullOpts.Store == nil {
		return nil, ErrStoreRequired
	}
	if fullOpts.Logger == nil {
		return nil, ErrLoggerRequired
	}

	// Set defaults
	if fullOpts.AdminHomeURL == "" {
		fullOpts.AdminHomeURL = "/admin"
	}
	if fullOpts.ShopAdminURL == "" {
		fullOpts.ShopAdminURL = "/admin/shop"
	}

	// Build controller routes ONCE at construction time (fixes #3 —
	// previously rebuilt on every request)
	uiConfig := shared.UiConfig{
		Store:            fullOpts.Store,
		Logger:           fullOpts.Logger,
		CustomerResolver: fullOpts.CustomerResolver,
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, layoutOpts struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			// If a custom layout is provided, use it directly (fixes #14 —
			// avoids computing the default layout when FuncLayout is set)
			if fullOpts.FuncLayout != nil {
				custom := fullOpts.FuncLayout(webpageTitle, webpageHtml, layoutOpts)
				if custom != "" {
					if w != nil {
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						_, _ = w.Write([]byte(custom))
						return ""
					}
					return custom
				}
			}

			webpage := shared.Layout(w, r, webpageTitle, webpageHtml, layoutOpts)

			// w may be nil when a subcontroller calls Layout() to get
			// the HTML string without writing to the response.
			if w != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(webpage))
				return ""
			}
			return webpage
		},
	}

	routes := buildControllerRoutes(uiConfig, fullOpts.FileManagerURL)

	handler := func(w http.ResponseWriter, r *http.Request) string {
		// Check authentication (fixes #2 — previously Routes() path
		// silently skipped the auth gate that admin.Handle() enforces)
		if fullOpts.AuthUserID != nil && fullOpts.AuthUserID(r) == "" {
			http.Redirect(w, r, fullOpts.AdminHomeURL, http.StatusSeeOther)
			return ""
		}

		// Inject config into context
		ctx := r.Context()
		ctx = contextWithValues(ctx, fullOpts, r.URL.Path)
		r = r.WithContext(ctx)

		controller := req.GetStringTrimmed(r, "controller")
		if controller == "" {
			controller = shared.CONTROLLER_HOME
		}

		fn, ok := routes[controller]
		if !ok {
			fn = routes[shared.CONTROLLER_HOME]
		}

		fn(w, r)
		return ""
	}

	shop := rtr.NewRoute().
		SetName("Admin > Shop").
		SetPath(fullOpts.ShopAdminURL).
		SetHTMLHandler(handler)

	shopCatchAll := rtr.NewRoute().
		SetName("Admin > Shop > Catchall").
		SetPath(fullOpts.ShopAdminURL + shared.CatchAll).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		shop,
		shopCatchAll,
	}, nil
}
