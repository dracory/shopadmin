# shopadmin

<img src="https://opengraph.githubassets.com/dracory/shopadmin" />

[![Tests Status](https://github.com/dracory/shopadmin/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/shopadmin/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/shopadmin)](https://goreportcard.com/report/github.com/dracory/shopadmin)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/shopadmin)](https://pkg.go.dev/github.com/dracory/shopadmin)

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Introduction

Standalone shop admin interface module for Go. Provides a ready-to-use
admin panel for managing products, categories, discounts, and orders,
built on top of [`github.com/dracory/shopstore`](https://github.com/dracory/shopstore).

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

Use `Routes()` to get `rtr.RouteInterface` slices for integration with
the host project's router:

```go
routes, err := shopadmin.Routes(app, shopadmin.AdminOptions{
    ShopAdminURL: links.Admin().Shop(),
    AdminHomeURL: links.Admin().Home(),
})
```

`app` must implement `shopadmin.RegistryInterface`:

```go
type RegistryInterface interface {
    GetShopStore() shopstore.StoreInterface
    GetCacheStore() cachestore.StoreInterface
    GetLogger() *slog.Logger
}
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

## Custom Layout

By default, shopadmin renders a bare-bones HTML page with Bootstrap and
Vue from CDN. To embed the admin inside your own layout (branding, menus,
etc.), provide `FuncLayout`:

```go
admin, _ := shopadmin.New(shopadmin.AdminOptions{
    Store:  store,
    Logger: logger,
    FuncLayout: func(title, body string, opts struct {
        Styles     []string
        StyleURLs  []string
        Scripts    []string
        ScriptURLs []string
    }) string {
        return myLayout(title, body, opts)
    },
})
```

The anonymous struct matches `cmsstore/admin` exactly, so you can reuse
your existing cmsstore layout function.

## Authentication

Provide an `AuthUserID` function to gate access. If it returns `""`, the
request is redirected to `AdminHomeURL`:

```go
admin, _ := shopadmin.New(shopadmin.AdminOptions{
    Store:      store,
    Logger:     logger,
    AuthUserID: func(r *http.Request) string {
        // return authenticated user ID, or ""
    },
})
```

## Project Structure

```
shopadmin/
├── shopadmin.go              # AdminOptions, New(), Handle()
├── routes.go                 # Routes() for router integration
├── registry.go               # RegistryInterface
├── types.go                  # CustomerResolverInterface (re-export)
├── errors.go                 # Sentinel errors
├── context.go                # Request context helpers
├── controllers.go            # buildControllerRoutes()
├── shared/                   # Shared UI infrastructure
│   ├── ui_config.go          # UiConfig, CustomerResolverInterface
│   ├── ui_base.go            # UiBase (embeds into controllers)
│   ├── ui_interface.go       # UiInterface
│   ├── layout.go             # Default layout
│   ├── flash.go              # Flash messages
│   ├── header.go             # Admin header
│   └── links.go              # URL builders
├── home/                     # Dashboard controller
├── product_manager/          # Product list controller
├── product_update/           # Product edit controller
├── product_delete/           # Product delete controller
├── category_manager/         # Category list controller
├── category_create/          # Category create controller
├── category_update/          # Category update controller
├── discount_manager/         # Discount list controller
├── order_manager/            # Order list controller
├── order_details/            # Order details controller
├── testutils/                # Test utilities
└── docs/
    └── proposal.md           # CustomerResolverInterface design doc
```

## Testing

```bash
go test ./...
```

Tests use an in-memory SQLite database via `modernc.org/sqlite` — no
external services required.

## Dependencies

- [`github.com/dracory/shopstore`](https://github.com/dracory/shopstore) — store interface
- [`github.com/dracory/cachestore`](https://github.com/dracory/cachestore) — flash messages (optional)
- [`github.com/dracory/rtr`](https://github.com/dracory/rtr) — router integration
- [`github.com/dracory/hb`](https://github.com/dracory/hb) — HTML builder
- [`github.com/dracory/bs`](https://github.com/dracory/bs) — Bootstrap components

**Not** dependent on `userstore` — customer resolution is via
`CustomerResolverInterface`.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

# Documentation

For more information, please refer to the [Documentation](./docs/proposal.md).
