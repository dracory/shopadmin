# Plan: Standalone `github.com/dracory/shopadmin` Module

**Task ID:** 2026-08-17-001
**Date:** 2026-08-17
**Status:** Approved v3 — implementation in progress

---

## 1. Summary

Port the in-repo `pkg/shopadmin` from `mechestudio.com` into a self-contained
Go module at `D:\PROJECTs\_modules_dracory\shopadmin` (module path
`github.com/dracory/shopadmin`). Sever all dependencies on
`project/internal/*` and `project/pkg/*`. Preserve the public API shape
(`AdminOptions`, `AdminInterface`, `New()`, `Routes()`) and the
folder-per-controller structure with embedded Vue.js SPA frontend.

The existing `pkg/shopadmin` in `mechestudio.com` stays **untouched**.

**This plan follows the established patterns from
`github.com/dracory/cmsstore/admin`** — the reference standalone admin
module in the dracory ecosystem.

---

## 2. Pattern Alignment with `cmsstore/admin`

The `cmsstore/admin` package solves the exact same problems (standalone
admin UI, pluggable layout, no host-project dependencies). I will follow
its conventions:

| Concern | `cmsstore/admin` pattern | shopadmin adaptation |
|---|---|---|
| **Entry point** | `AdminOptions` struct, `New(options) (*admin, error)`, `Handle(w, r)` | Same — already matches in-repo shopadmin |
| **Store dependency** | `Store cmsstore.StoreInterface` field directly in `AdminOptions` (not a registry) | `Store shopstore.StoreInterface` + `UserStore userstore.StoreInterface` directly in `AdminOptions` |
| **Layout** | `FuncLayout` optional function field; default `shared.Layout()` renders a bare HTML page with Bootstrap CDN | Same pattern: `FuncLayout` optional field; default `shared.Layout()` renders bare HTML with Bootstrap + Vue CDN |
| **URLs** | `shared.URL()` and `shared.URLR()` build URLs from request context (`KeyEndpoint`) | Same pattern — endpoint stored in request context, `shared.URLR()` builds controller URLs |
| **Breadcrumbs** | `shared.Breadcrumb` struct + `shared.Breadcrumbs()` function in `shared/` package | Same — move `Breadcrumb` and `Breadcrumbs()` into `shared/` |
| **Config via context** | `KeyAdminHomeURL`, `KeyMediaManagerURL`, etc. injected into request context in `Handle()` | Same — `KeyAdminHomeURL`, `KeyFileManagerURL`, `KeyShopAdminURL` injected in `Handle()` |
| **Controller deps** | `UiConfig` struct passed to `UI(config)` factory; controllers receive `UiInterface` | `UiConfig` struct with `Store`, `UserStore`, `Layout`, `Logger`; controllers receive `UiInterface` |
| **Tests** | `testutils.InitStore(":memory:")` in test files; `newAdminForTest()` helper | Same — `testutils.InitStore(":memory:")` + `newAdminForTest()` helper |
| **No `dashboard` dependency** | `cmsstore/admin` uses `hb.NewWebpage()` for its default layout, NOT `github.com/dracory/dashboard` | **Drop `github.com/dracory/dashboard` from the plan.** Use `hb.NewWebpage()` like cmsstore does. |

### Key change from plan v1

**`github.com/dracory/dashboard` is no longer needed.** The `cmsstore/admin`
module uses `hb.NewWebpage()` to build a complete HTML page with Bootstrap
CDN — no dashboard framework dependency. I'll do the same. This removes the
only flagged new dependency.

---

## 3. Current Dependency Analysis

The in-repo `pkg/shopadmin` imports four `project/internal/*` packages and
one `project/pkg/*` package:

| Dependency | What shopadmin uses from it | Replacement (following cmsstore pattern) |
|---|---|---|
| `project/internal/app` | `app.AppInterface` — only 4 methods called: `GetShopStore()`, `GetUserStore()`, `GetCacheStore()`, `GetLogger()` | **Drop the registry.** Put `Store`, `UserStore`, `CacheStore`, `Logger` directly in `AdminOptions` (like cmsstore puts `Store` in `AdminOptions`) |
| `project/internal/links` | `ADMIN_SHOP` const, `CATCHALL` const, `Admin().Home()`, `Admin().Shop(params)` | Constants + `shared.URL()`/`shared.URLR()` in `shared/` package |
| `project/internal/layouts` | `NewAdminLayout(app, r, Options).ToHTML()`, `Breadcrumbs([]Breadcrumb)`, `Breadcrumb` struct, `Options` struct | `FuncLayout` optional field (like cmsstore); `shared.Breadcrumbs()` + `shared.Breadcrumb` in `shared/` |
| `project/internal/helpers` | `GetAuthUser(r)` (nil-check), `ToFlashError(cacheStore, w, r, msg, url, seconds)` | `AuthUserID` callback (already in `AdminOptions`); flash via `CacheStore` + redirect (internal helper) |
| `project/internal/testutils` | `testutils.Setup(...)` in tests only | `testutils.InitStore(":memory:")` pattern from cmsstore |
| `project/pkg/dashboard` | `dashboard.DashboardInterface`, `dashboard.New()` | **Not needed** — use `hb.NewWebpage()` like cmsstore |

### External `github.com/dracory/*` dependencies (all already in the ecosystem)

- `github.com/dracory/shopstore` — store interfaces, entity types, query builders
- `github.com/dracory/userstore` — `UserInterface`, `NewUserQuery()` (order controllers)
- `github.com/dracory/cachestore` — `StoreInterface` (flash messages)
- `github.com/dracory/rtr` — route types (for `Routes()` return type)
- `github.com/dracory/req` — request parameter helpers
- `github.com/dracory/hb` — HTML builder (includes `NewWebpage()` for default layout)
- `github.com/dracory/bs` — Bootstrap components (modals, nav tabs)
- `github.com/dracory/api` — JSON API responses
- `github.com/dracory/cdn` — CDN URLs for Vue, htmx, SweetAlert, jQuery, Notiflix
- `github.com/dracory/neat` — `SortDesc` constant
- `github.com/dracory/uid` — `HumanUid()`
- `github.com/dracory/test` — test helpers (used by cmsstore tests)

**No new dependencies.** Every dependency is already used by `cmsstore/admin`
or `shopstore`. The `github.com/dracory/dashboard` dependency from plan v1 is
dropped.

### Other external dependencies

- `github.com/samber/lo` — `FirstOr`, `IfF`, `HasKey`
- `github.com/spf13/cast` — `ToString`

---

## 4. Design Decisions

### 4.1 AdminOptions (follows cmsstore pattern)

```go
type AdminOptions struct {
    // Store is the shopstore.StoreInterface (required)
    Store shopstore.StoreInterface

    // UserStore is optional — needed only for order controllers
    // to resolve customer names/emails
    UserStore userstore.StoreInterface

    // CacheStore is optional — used for flash messages.
    // If nil, flash errors fall back to inline HTML rendering.
    CacheStore cachestore.StoreInterface

    // Logger is required (matches cmsstore requirement)
    Logger *slog.Logger

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
```

**Why this is better than plan v1's `RegistryInterface`:**
- Matches the established `cmsstore/admin` pattern exactly
- Simpler — no new interface type to document and maintain
- Consumers pass concrete stores directly, not a registry object
- Each store is independently optional (e.g., `UserStore` is nil for
  projects that don't need order customer resolution)
- `FuncLayout` uses anonymous struct (not named `LayoutOptions` type) to
  match cmsstore/admin exactly — consumers can reuse their existing
  cmsstore layout function for shopadmin without adaptation

**Migration for existing consumers:**
```go
// Before (in-repo):
admin, err := shopadmin.New(shopadmin.AdminOptions{
    Registry: app,  // app.AppInterface
})

// After (standalone):
admin, err := shopadmin.New(shopadmin.AdminOptions{
    Store:      app.GetShopStore(),
    UserStore:  app.GetUserStore(),
    CacheStore: app.GetCacheStore(),
    Logger:     app.GetLogger(),
})
```

This is a **breaking change at the call site** — consumers must change
from passing `Registry: app` to passing individual stores. This is
unavoidable to sever the `project/internal/app` dependency, and it
matches how `cmsstore/admin` works. The migration is mechanical and
obvious. This breaking change is **acknowledged as a deviation from the
task's "keep signatures stable" requirement** — it is unavoidable because
`app.AppInterface` is a `project/internal/*` type that cannot be imported
by a standalone module.

### 4.2 Layout (follows cmsstore pattern exactly)

`FuncLayout` uses an **anonymous struct** for layout options — matching
`cmsstore/admin` exactly so consumers can reuse their existing cmsstore
layout function:

```go
FuncLayout func(title string, body string, options struct {
    Styles     []string
    StyleURLs  []string
    Scripts    []string
    ScriptURLs []string
}) string
```

The default layout (when `FuncLayout` is nil) uses `hb.NewWebpage()` with
Bootstrap + Vue CDN — exactly like `cmsstore/admin/shared/layout.go`.

### 4.3 URLs (follows cmsstore pattern)

```go
// shared/consts.go
const KeyEndpoint      = "endpoint"
const KeyAdminHomeURL  = "admin_home_url"
const KeyShopAdminURL  = "shop_admin_url"
const KeyFileManagerURL = "file_manager_url"

const PathHome           = "home"
const PathProducts       = "products"
const PathProductUpdate  = "product_update"
const PathProductDelete  = "product_delete"
const PathCategories     = "categories"
const PathCategoryCreate = "category_create"
const PathCategoryUpdate = "category_update"
const PathDiscounts      = "discounts"
const PathOrders         = "orders"
const PathOrderDetails   = "order_details"

const CatchAll = "/*"
```

```go
// shared/URL.go
func URL(endpoint, controller string, params map[string]string) string
func URLR(r *http.Request, controller string, params map[string]string) string
```

`URLR` reads the endpoint from request context (injected by `Handle()`),
same as `cmsstore/admin`. The `controller` value is placed in a **copy**
of the `params` map under the `"controller"` key (fixes pre-existing bug
#10 where the caller's map was mutated).

**Base URL from context, not hardcoded:** The in-repo code hardcodes
`shared.NewLinks("/admin/shop")` in every subcontroller. In the
standalone module, the base URL is read from request context via
`shared.ShopAdminURL(r)` (injected in `Handle()`). The `shared.Links`
helper struct is either constructed with the context-derived URL or
replaced entirely by `shared.URLR()` calls.

### 4.4 Config via context + routes built once (follows cmsstore pattern)

**Routes built once in `New()`, not per-request.** The in-repo `Handle()`
calls `Routes()` on every HTTP request, creating new route objects each
time. The standalone module follows cmsstore's pattern: the route map is
built once in `New()` and stored on the `admin` struct. `Handle()` just
looks up the map — no allocation per request.

In `Handle()`:
```go
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
    ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
    ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.adminHomeURL)
    ctx = context.WithValue(ctx, shared.KeyShopAdminURL, a.shopAdminURL)
    ctx = context.WithValue(ctx, shared.KeyFileManagerURL, a.fileManagerURL)
    r = r.WithContext(ctx)

    // Map-based lookup (like cmsstore), not HasPrefix matching
    controller := req.GetStringTrimmed(r, "controller")
    handler := a.getRoute(controller) // looks up pre-built map
    handler(w, r)
}
```

This also fixes pre-existing bug #13: the in-repo `HasPrefix` matching
could match `/admin/shopanything` to the `/admin/shop` route. The
map-based lookup is exact-match on the controller name.

### 4.5 UiConfig / UiInterface (follows cmsstore pattern)

```go
// shared/ui_config.go
type UiConfig struct {
    Store      shopstore.StoreInterface
    UserStore  userstore.StoreInterface
    CacheStore cachestore.StoreInterface
    Logger     *slog.Logger
    Layout     func(w http.ResponseWriter, r *http.Request, title, body string, options struct {
        Styles     []string
        StyleURLs  []string
        Scripts    []string
        ScriptURLs []string
    }) string
}
```

Uses anonymous struct for layout options — matches `cmsstore/admin` exactly.

Each subcontroller package (e.g., `home/`, `product_manager/`) exposes:
```go
func UI(config shared.UiConfig) UiInterface
```

Controllers receive `UiInterface` and call `ui.Store()`, `ui.Logger()`,
`ui.Layout(...)`. This is the exact pattern from
`cmsstore/admin/blocks/ui.go`.

### 4.6 Auth & Flash

**Auth:** Subcontrollers check `AuthUserID` from request context. In
`Handle()`, if `AdminOptions.AuthUserID` is set, the result is stored in
context. Subcontrollers check whether it's empty.

**Flash:** Internal `toFlashError()` helper. If `CacheStore` is
available, stores the message and redirects to a flash URL. If not,
renders the error inline as HTML. The flash URL is derived from
`AdminHomeURL + "/flash"` (configurable). This is simpler than plan v1's
separate `FlashURL` field — and matches the cmsstore philosophy of
keeping options minimal.

**Nil-safety (fixes pre-existing bugs #11, #14):**
- `toFlashError()` nil-checks the `http.ResponseWriter` before calling
  `http.Redirect()` — several in-repo controllers pass `nil` as the
  writer, which would panic.
- `shared.Header()` nil-checks the logger before calling `logger.Error()`
  — even though `Logger` is required in `AdminOptions`, the public
  `Header()` function should be defensive.

### 4.7 Routes() function — signature preserved

The in-repo `Routes()` returns `[]rtr.RouteInterface` for integration
with the host project's router. The task requires keeping this signature
stable. The in-repo signature is:

```go
func Routes(app app.AppInterface, opts ...AdminOptions) ([]rtr.RouteInterface, error)
```

Since `app.AppInterface` is a `project/internal/*` type that cannot be
imported, the first parameter changes to a structurally-compatible
interface defined in shopadmin:

```go
// RegistryInterface satisfies the same structural interface as
// project/internal/app.AppInterface for the methods Routes() needs.
// Any app.AppInterface implementation satisfies this.
type RegistryInterface interface {
    GetShopStore() shopstore.StoreInterface
    GetUserStore() userstore.StoreInterface
    GetCacheStore() cachestore.StoreInterface
    GetLogger() *slog.Logger
}

func Routes(registry RegistryInterface, opts ...AdminOptions) ([]rtr.RouteInterface, error)
```

**This preserves the call-site signature** — consumers who currently
call `shopadmin.Routes(app, opts)` can continue to do so unchanged,
because `app.AppInterface` structurally satisfies `RegistryInterface`.
The `opts` variadic parameter is also preserved.

Internally, `Routes()` builds the `UiConfig` from `registry` and `opts`,
and dispatches to subcontrollers via the `?controller=` query parameter,
same as before.

### 4.8 New() defaults + test harness

**`New()` defaults (fixes pre-existing bug #12):** The in-repo `New()`
only sets a default for `ShopAdminURL`, not `AdminHomeURL`. The
standalone `New()` sets both defaults:
```go
if opts.AdminHomeURL == "" {
    opts.AdminHomeURL = "/admin"
}
if opts.ShopAdminURL == "" {
    opts.ShopAdminURL = "/admin/shop"
}
```

**Test harness (follows cmsstore pattern):**
`testutils/utils.go` at module root (exported, like cmsstore):
```go
func InitStore(filepath string) (shopstore.StoreInterface, error)
```

Test files use `testutils.InitStore(":memory:")` + a local
`newAdminForTest()` helper — exactly like `cmsstore/admin/admin_test.go`.

Per `AGENTS.md`: "Test cases are difficult to maintain, use only test
functions" — we write test **functions** covering:
- `New()` validation (nil store, nil logger, default URLs)
- `Routes()` registration (returns 2 routes, correct paths)
- One subcontroller smoke test (home controller renders without panic)

### 4.9 Pre-existing bug fixes during port

**Fix (crash/security risks):**
- **#6 — unchecked type assertions in `handleSaveMedia`:** Lines that do
  `item["url"].(string)` will panic if the field is missing or wrong type.
  Fix: use safe assertions `if url, ok := item["url"].(string); ok { ... }`.
- **#7 — error message injection into JSON in `handleUploadMedia`:**
  `err.Error()` is concatenated into a JSON string literal, producing
  malformed JSON if the error contains quotes. Fix: use `json.Marshal` for
  the response body.

**Defer (logic bugs — fixing changes behavior, needs separate testing):**
- **#8 — `OrderCount` uses unfiltered query:** The total count doesn't
  match the filtered results, breaking pagination. Deferred because fixing
  it changes the pagination behavior and needs dedicated testing.
- **#9 — `OrderLineItemList` loads all line items:** Loads every line item
  in the DB then filters in Go. Should use `SetOrderID()` on the query.
  Deferred because it's a performance fix, not a crash, and changing the
  query could surface other issues.

**Fixed by design (addressed in §4.3, §4.4, §4.6):**
- #10 — map mutation in `shared/links.go` → fixed by copying params map
- #11 — nil logger in `Header()` → fixed by nil-check
- #12 — `AdminHomeURL` default → fixed in `New()`
- #13 — `HasPrefix` route matching → fixed by map-based lookup
- #14 — nil writer in `ToFlashError` → fixed by nil-check

---

## 5. File Structure (New Module)

```
D:\PROJECTs\_modules_dracory\shopadmin\
├── go.mod                          # module github.com/dracory/shopadmin, go 1.26
├── go.sum                          # via go mod tidy
├── LICENSE                         # carried over (already exists)
├── README.md                       # new, mirrors in-repo README + migration note
├── plan.md                         # this file
├── shopadmin.go                    # AdminOptions, AdminInterface, New(), Handle()
├── registry.go                     # RegistryInterface (for Routes() signature)
├── errors.go                       # error definitions
├── routes.go                       # Routes() dispatcher
├── shopadmin_test.go               # New() + Routes() tests
├── testutils/
│   └── utils.go                    # InitStore() for tests
├── shared/
│   ├── consts.go                   # context keys, path constants, CatchAll
│   ├── URL.go                      # URL(), URLR(), query()
│   ├── Endpoint.go                 # Endpoint(r)
│   ├── admin_home_url.go           # AdminHomeURL(r)
│   ├── shop_admin_url.go           # ShopAdminURL(r)
│   ├── file_manager_url.go         # FileManagerURL(r)
│   ├── Breadcrumb.go               # Breadcrumb struct
│   ├── Breadcrumbs.go              # Breadcrumbs() function
│   ├── ui_config.go                # UiConfig struct
│   ├── ui_interface.go             # UiInterface interface
│   ├── layout.go                   # Layout() default renderer (hb.NewWebpage)
│   ├── constants.go                # CONTROLLER_* constants
│   ├── links.go                    # link helpers (self-contained)
│   ├── header.go                   # Header() nav component
│   └── flash.go                    # toFlashError() helper
├── home/
│   ├── ui.go                       # UI(config) factory
│   ├── home_controller.go
│   ├── home.html                   # embedded
│   ├── home.js                     # embedded
│   └── home_test.go                # smoke test
├── product_manager/
│   ├── ui.go
│   ├── constants.go
│   ├── product_manager_controller.go
│   ├── product_manager_page.go
│   ├── handle_products_fetch_ajax.go
│   ├── handle_product_create_ajax.go
│   ├── handle_product_delete_ajax.go
│   ├── handle_product_delete_selected_ajax.go
│   ├── products.html               # embedded
│   └── products.js                 # embedded
├── product_update/
│   ├── ui.go
│   ├── product_update_controller.go
│   ├── product_details_component.go
│   ├── product_media_component.go
│   ├── product_metadata_component.go
│   ├── product_tags_component.go
│   ├── types.go
│   ├── details.html, details.js    # embedded
│   ├── media.html, media.js        # embedded
│   ├── metadata.html, metadata.js  # embedded
│   └── tags.html, tags.js          # embedded
├── product_delete/
│   ├── ui.go
│   └── product_delete_controller.go
├── category_manager/
│   ├── ui.go
│   ├── constants.go
│   ├── category_manager_controller.go
│   ├── category_manager_page.go
│   ├── handle_categories_load_ajax.go
│   ├── handle_category_delete_ajax.go
│   ├── handle_category_delete_selected_ajax.go
│   ├── categories.html             # embedded
│   └── categories.js               # embedded
├── category_create/
│   ├── ui.go
│   └── category_create_controller.go
├── category_update/
│   ├── ui.go
│   └── category_update_controller.go
├── discount_manager/
│   ├── ui.go
│   ├── constants.go
│   ├── discount_manager_controller.go
│   ├── discount_manager_page.go
│   ├── handle_discounts_load_ajax.go
│   ├── handle_discount_delete_ajax.go
│   ├── handle_discount_delete_selected_ajax.go
│   ├── discounts.html              # embedded
│   └── discounts.js                # embedded
├── order_manager/
│   ├── ui.go
│   ├── constants.go
│   ├── order_manager_controller.go
│   ├── order_manager_page.go
│   ├── handle_orders_load_ajax.go
│   ├── orders.html                 # embedded
│   └── orders.js                   # embedded
└── order_details/
    ├── ui.go
    ├── order_details_controller.go
    ├── order_details_page.go
    ├── handle_order_details_load_ajax.go
    ├── order_details.html          # embedded
    └── order_details.js            # embedded
```

---

## 6. Public API (Before → After)

### Before (in-repo)
```go
import "project/pkg/shopadmin"

admin, err := shopadmin.New(shopadmin.AdminOptions{
    Registry:      app,              // app.AppInterface
    AdminHomeURL:  "/admin",
    ShopAdminURL:  "/admin/shop",
    AuthUserID:    func(r *http.Request) string { ... },
    FileManagerURL: "/admin/files",
})
```

### After (standalone)
```go
import "github.com/dracory/shopadmin"

admin, err := shopadmin.New(shopadmin.AdminOptions{
    Store:        app.GetShopStore(),
    UserStore:    app.GetUserStore(),
    CacheStore:   app.GetCacheStore(),
    Logger:       app.GetLogger(),
    AdminHomeURL: "/admin",
    ShopAdminURL: "/admin/shop",
    AuthUserID:   func(r *http.Request) string { ... },
    FileManagerURL: "/admin/files",
    FuncLayout:   myLayoutFunc,  // optional
})
```

**Breaking changes at call site:** Yes — `Registry` field is replaced by
individual `Store`, `UserStore`, `CacheStore`, `Logger` fields. This is
unavoidable and matches the `cmsstore/admin` convention. Migration is
mechanical.

**New optional fields:**
- `FuncLayout` — custom layout rendering; defaults to built-in
  `hb.NewWebpage()`-based layout (Bootstrap + Vue CDN)
- `UserStore` — optional; needed only for order customer resolution
- `CacheStore` — optional; needed only for flash messages

---

## 7. Implementation Steps

### Step 1: Initialize module
- Write `go.mod` (module `github.com/dracory/shopadmin`, go 1.26)
- Run `go mod tidy` after all files are in place

### Step 2: Port `shared/` package (follows cmsstore pattern)
- `consts.go` — context keys, path constants, `CatchAll`
- `URL.go` — `URL()`, `URLR()`, `query()` (adapted from cmsstore)
- `Endpoint.go` — `Endpoint(r)`
- `admin_home_url.go`, `shop_admin_url.go`, `file_manager_url.go`
- `Breadcrumb.go`, `Breadcrumbs.go` — moved from `project/internal/layouts`
- `ui_config.go`, `ui_interface.go` — new, follows cmsstore pattern
- `layout.go` — default layout using `hb.NewWebpage()` (adapted from
  `cmsstore/admin/shared/layout.go`)
- `constants.go` — `CONTROLLER_*` constants (direct copy)
- `links.go` — link helpers using `shared.URL()` (rewrite)
- `header.go` — nav header (port, use `shared.URLR()`)
- `flash.go` — `toFlashError()` helper

### Step 3: Port core files
- `errors.go` — direct copy, update error messages
- `shopadmin.go` — port, replace `Registry` with individual store fields,
  inject config into request context in `Handle()` (like cmsstore)
- `routes.go` — port, replace `app` param with `opts AdminOptions`,
  build `UiConfig` from opts, dispatch to subcontrollers

### Step 4: Port all 10 subcontrollers
For each subcontroller:
- Create `ui.go` with `UI(config shared.UiConfig) UiInterface` factory
  (follows `cmsstore/admin/blocks/ui.go` pattern)
- Copy `.go` files and embedded `.html`/`.js` files
- Replace `project/internal/app` → use `ui.Store()`, `ui.UserStore()`,
  `ui.CacheStore()`, `ui.Logger()` via `UiInterface`
- Replace `project/internal/helpers` → use `AuthUserID` from context +
  internal `toFlashError()`
- Replace `project/internal/layouts` → use `ui.Layout()` +
  `shared.Breadcrumbs()`
- Replace `project/internal/links` → use `shared.URLR()`
- Replace `project/pkg/shopadmin/shared` → `github.com/dracory/shopadmin/shared`

### Step 5: Write testutils + tests
- `testutils/utils.go` — `InitStore(":memory:")` (follows cmsstore)
- `shopadmin_test.go` — `TestNew_RequiresStore`, `TestNew_RequiresLogger`,
  `TestNew_DefaultURLs`, `TestRoutes_ReturnsRoutes`
- `home/home_test.go` — smoke test: render home page with test store

### Step 6: Write README.md
- Mirror in-repo README structure
- Update import paths to `github.com/dracory/shopadmin`
- Add migration note showing the `Registry` → individual stores change
- Document `FuncLayout`, `UserStore`, `CacheStore` optional fields
- Document `UiConfig`/`UiInterface` pattern for extension

### Step 7: Build & test
- `go mod tidy`
- `go build ./...`
- `go test ./...`
- Fix any issues

### Step 8: Verify in-repo untouched
- Confirm `D:\PROJECTs\mechestudio.com\pkg\shopadmin\` has no git changes

---

## 8. What Stays the Same

- Public API shape: `AdminOptions`, `AdminInterface`, `New()`, `Routes()`
- Folder-per-controller structure
- Embedded Vue.js SPA frontend (all `.html`/`.js` files copied verbatim)
- All AJAX handler logic (product/category/discount/order CRUD)
- Route paths (`/admin/shop`, `/admin/shop/*`)
- Controller dispatch via `?controller=` query parameter
- All external `github.com/dracory/*` dependencies already in use

---

## 9. What Changes

| Area | Before | After | Impact on consumers |
|---|---|---|---|
| `AdminOptions.Registry` | `app.AppInterface` | Removed; replaced by `Store`, `UserStore`, `CacheStore`, `Logger` | **Breaking** — mechanical migration: `Registry: app` → `Store: app.GetShopStore(), ...` |
| `Routes()` first param | `app.AppInterface` | `RegistryInterface` (structurally compatible) | **None** — `app.AppInterface` satisfies `RegistryInterface`, call site unchanged |
| Subcontroller constructors | `NewHomeController(app, url)` | `UI(config).HomeController` via `UiInterface` | None — internal to module |
| Layout rendering | Hardcoded `layouts.NewAdminLayout()` | `FuncLayout` optional (anonymous struct, matches cmsstore) + default `hb.NewWebpage()` | None if default is acceptable; optional custom renderer |
| Auth check | `helpers.GetAuthUser(r) != nil` | `AuthUserID(r) != ""` from context | None — `AuthUserID` already in `AdminOptions` |
| Flash messages | `helpers.ToFlashError()` | Internal `toFlashError()` via `CacheStore` (nil-safe) | None — `CacheStore` is now an explicit field |
| Link constants | `links.ADMIN_SHOP`, `links.CATCHALL` | `shared` constants | None — internal |
| Link base URL | Hardcoded `"/admin/shop"` | From request context (`shared.ShopAdminURL(r)`) | None — internal |
| Route matching | `strings.HasPrefix` (ambiguous) | Map-based exact match (like cmsstore) | None — internal, fixes bug #13 |
| Route building | Per-request in `Handle()` | Once in `New()`, stored on struct | None — internal, fixes perf issue #4 |
| Test setup | `testutils.Setup(...)` | `testutils.InitStore(":memory:")` | None — test-only |
| Config passing | `app` passed to every controller | Config in request context (like cmsstore) | None — internal |

---

## 10. Risks & Mitigations

| Risk | Mitigation |
|---|---|
| `Registry` → individual stores is a breaking change | Mechanical migration; documented in README; matches cmsstore convention |
| Default layout looks different from host project | `FuncLayout` is pluggable; consumers wire their own. Default is a clean fallback. |
| Flash without `CacheStore` | Falls back to inline HTML error rendering. Graceful degradation. |
| Order controllers need `UserStore` | `UserStore` is optional in `AdminOptions`; order controllers nil-check before use (they already do this). |

---

## 11. Will NOT Do

- Will **not** delete or modify `D:\PROJECTs\mechestudio.com\pkg\shopadmin\`
- Will **not** push to GitHub (await MD/GOD approval)
- Will **not** introduce new dependencies (all deps already in cmsstore/shopstore)
- Will **not** change the Vue.js frontend behavior
- Will **not** change route paths or controller dispatch logic
- Will **not** add table-driven test suites with complex fixtures (per AGENTS.md: "use only test functions")

---

## 12. Definition of Done

- [ ] All subcontrollers ported; no `project/*` imports remain
- [ ] `go build ./...` passes in the new module
- [ ] `go test ./...` passes in the new module
- [ ] `README.md` updated with new import path and migration note
- [ ] Existing `pkg/shopadmin` in `mechestudio.com` untouched
- [ ] Report delivered to MD

---

## 13. Items Flagged for MD Review (resolved)

1. **`Registry` field removed, replaced by individual stores** — This is
   a breaking change at the call site. Consumers must change from
   `Registry: app` to `Store: app.GetShopStore(), UserStore: app.GetUserStore(), ...`.
   This matches the `cmsstore/admin` convention and is unavoidable to
   sever the `project/internal/app` dependency. **Acknowledged as a
   deviation from the task's "keep signatures stable" requirement.**
   **MD decision: approved.**

2. **`Routes()` signature preserved** — First param changes from
   `app.AppInterface` to `RegistryInterface` (structurally compatible).
   Call site unchanged: `shopadmin.Routes(app, opts)` still works.
   **MD decision: approved (no change needed at call site).**

3. **`FuncLayout` uses anonymous struct** — Matches `cmsstore/admin`
   exactly so consumers can reuse their existing cmsstore layout function.
   No new dependencies. Default uses `hb.NewWebpage()` with Bootstrap +
   Vue CDN. **MD decision: approved.**

4. **No `github.com/dracory/dashboard` dependency** — Dropped from plan
   v1. Using `hb.NewWebpage()` instead, exactly like `cmsstore/admin`.
   **MD decision: approved.**

5. **Pre-existing bug fixes** — Fixing crash/security risks (#6 unchecked
   type assertions, #7 JSON injection). Deferring logic bugs (#8 order
   count, #9 line item query) to avoid behavior changes without dedicated
   testing. **MD decision: approved.**

6. **Subcontroller constructor pattern change** — Switching from
   `NewHomeController(app, url)` to `UI(config).HomeController` via
   `UiInterface`. Internal to the module, but exported types change.
   Consumers should only call `shopadmin.New()` and `shopadmin.Routes()`.
   **MD decision: approved.**
