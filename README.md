# shopadmin

<img src="https://opengraph.githubassets.com/dracory/shopadmin" />

[![Tests Status](https://github.com/dracory/shopadmin/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/shopadmin/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/shopadmin)](https://goreportcard.com/report/github.com/dracory/shopadmin)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/shopadmin)](https://pkg.go.dev/github.com/dracory/shopadmin)

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Introduction

Admin interface for [`github.com/dracory/shopstore`](https://github.com/dracory/shopstore).
Provides a ready-to-use admin panel for managing products, categories,
discounts, and orders.

Modeled after [`github.com/dracory/cmsstore/admin`](https://github.com/dracory/cmsstore)
— same folder-per-controller pattern, same `UiConfig`/`UiBase` conventions.

## Features

- **Product management** — create, update, delete, list with AJAX
- **Category management** — create, update, delete, list with AJAX
- **Discount management** — create, delete, list with AJAX
- **Order management** — list orders, view order details
- **Media management** — upload, reorder, delete product images
- **Metadata & tags** — per-product metadata and tag editing
- **Customer resolution** — pluggable via `CustomerResolverInterface`
  (no dependency on any specific user/auth package)
- **Custom layouts** — bring your own layout via `FuncLayout`
- **Bootstrap + Vue CDN** — default UI works out of the box

## Installation

```bash
go get github.com/dracory/shopadmin
```

## Quick Start

```go
package main

import (
    "log/slog"
    "net/http"
    "os"

    "github.com/dracory/shopadmin"
    "github.com/dracory/shopstore"
)

func main() {
    store, err := shopstore.NewStore(shopstore.NewStoreOptions{
        DB:                 yourDB,
        ProductTableName:   "shop_product",
        CategoryTableName:  "shop_category",
        // ... other table names
        AutomigrateEnabled: true,
    })
    if err != nil {
        log.Fatal(err)
    }

    admin, err := shopadmin.New(shopadmin.AdminOptions{
        Store:       store,
        Logger:      slog.New(slog.NewTextHandler(os.Stderr, nil)),
        AdminHomeURL: "/admin",
        ShopAdminURL: "/admin/shop",
    })
    if err != nil {
        log.Fatal(err)
    }

    http.Handle("/admin/shop", admin)
    http.ListenAndServe(":8080", nil)
}
```

## Integration with a Router

`shopadmin.AdminInterface` exposes `Handle(w, r)`, which is an
`http.HandlerFunc`-compatible method. Wire it into any router that
accepts standard `http.Handler`:

```go
// stdlib
mux.Handle("/admin/shop", http.HandlerFunc(admin.Handle))

// github.com/dracory/rtr
route := rtr.NewRoute().
    SetName("Admin > Shop").
    SetPath("/admin/shop").
    SetHTMLHandler(admin.Handle)
```

## Customer Resolution

Shopadmin does **not** depend on `userstore` or any specific auth package.
Instead, order controllers resolve customer names/emails via an optional
`CustomerResolverInterface`:

```go
type CustomerResolverInterface interface {
    FindByID(ctx context.Context, customerID string) (name, email string)
    SearchIDs(ctx context.Context, name, email string) ([]string, error)
}
```

Provide an implementation at construction time:

```go
admin, _ := shopadmin.New(shopadmin.AdminOptions{
    Store:  store,
    Logger: logger,
    CustomerResolver: &myCustomerResolver{userStore: app.GetUserStore()},
})
```

If `CustomerResolver` is nil, customer fields stay empty and customer
filtering is disabled — no panic, no error.

See [`docs/proposal.md`](docs/proposal.md) for the full design rationale.

## Custom Sections for Product View and Update Pages

`shopadmin` allows host applications to inject custom content and dynamic UI components into product view and product update pages via `ProductViewSections` and `ProductUpdateSections` on `AdminOptions`.

### Supported Positions

- `shopadmin.PositionTop` (`"top"`): Rendered at the top of the page content.
- `shopadmin.PositionBeforeForm` (`"before_form"`): Rendered before the main update form / card.
- `shopadmin.PositionAfterForm` (`"after_form"`): Rendered after the main update form / card.
- `shopadmin.PositionSidebar` (`"sidebar"`): Rendered in a sidebar column alongside the product card.
- `shopadmin.PositionBottom` (`"bottom"`): Rendered at the bottom of the page container.

### Registering Custom Sections

You can register section renderers that receive the current `*http.Request` and `shopstore.ProductInterface`:

```go
admin, _ := shopadmin.New(shopadmin.AdminOptions{
    Store:  store,
    Logger: logger,
    ProductViewSections: []shopadmin.ProductSectionRenderer{
        // Simple HTML banner section
        func(r *http.Request, p shopstore.ProductInterface) (shopadmin.ProductSection, error) {
            return shopadmin.ProductSection{
                ID:       "inventory-notice",
                Position: shopadmin.PositionTop,
                HTML:     `<div class="alert alert-info">Viewing Product ID: ` + p.GetID() + `</div>`,
            }, nil
        },
    },
    ProductUpdateSections: []shopadmin.ProductSectionRenderer{
        // Advanced interactive widget section with CSS, JS, and third-party assets
        func(r *http.Request, p shopstore.ProductInterface) (shopadmin.ProductSection, error) {
            return shopadmin.ProductSection{
                ID:       "analytics-widget",
                Position: shopadmin.PositionSidebar,
                HTML:     `<div id="product-chart-card" class="card"><div class="card-body"><canvas id="salesChart"></canvas></div></div>`,
                CSS:      `#product-chart-card { margin-top: 10px; }`,
                JS:       `new Chart(document.getElementById('salesChart'), { type: 'line', data: { ... } });`,
                CSSURLs:  []string{"https://cdn.jsdelivr.net/npm/chart.js/dist/chart.min.css"},
                JSURLs:   []string{"https://cdn.jsdelivr.net/npm/chart.js/dist/chart.umd.min.js"},
            }, nil
        },
    },
})
```

## Custom Layout

By default, shopadmin renders a bare-bones HTML page with Bootstrap and
Vue from CDN. To embed the admin inside your own layout (branding, menus,
etc.), provide `FuncLayout`:

```go
admin, _ := shopadmin.New(shopadmin.AdminOptions{
    Store:  store,
    Logger: logger,
    FuncLayout: func(w http.ResponseWriter, r *http.Request, title, body string, opts struct {
        Styles     []string
        StyleURLs  []string
        Scripts    []string
        ScriptURLs []string
    }) string {
        return myLayout(w, r, title, body, opts)
    },
})
```

The anonymous struct matches `cmsstore/admin` exactly, so you can reuse
your existing cmsstore layout function.

## Testing

```bash
go test ./...
```

Tests use an in-memory SQLite database via `modernc.org/sqlite` — no
external services required.
