# Proposal: Remove `userstore` Dependency — Replace with CustomerResolverInterface

**Date:** 2026-08-17
**Status:** Draft — awaiting MD/GOD approval
**Module:** `github.com/dracory/shopadmin`

---

## 1. Problem

The in-repo `pkg/shopadmin` imports `github.com/dracory/userstore` for exactly
3 methods, used in exactly 2 files:

| Method | Used In | Purpose |
|--------|---------|---------|
| `UserFindByID(ctx, id)` | `order_details/handle_order_details_load_ajax.go:58-63` | Resolve customer name/email for a single order |
| `UserList(ctx, query)` | `order_manager/handle_orders_load_ajax.go:96-117` | Filter orders by customer name/email substring |
| `NewUserQuery()` | `order_manager/handle_orders_load_ajax.go:99` | Constructor for the query |

Three fields accessed on the returned user:
- `GetID()`
- `GetFirstName()` + `GetLastName()`
- `GetEmail()`

Shopadmin doesn't care about **users** — it cares about **customers**. The
customer data source (userstore today, maybe a CRM tomorrow) is an
implementation detail of the host project. Forcing `userstore` as a
dependency couples shopadmin to a specific auth/user-management package it
has no business caring about.

---

## 2. Solution

Replace the `userstore` dependency with a `CustomerResolverInterface`
defined in shopadmin. The host project provides an implementation at
construction time.

```go
// CustomerResolverInterface resolves customer data for order views.
// The host project provides an implementation — shopadmin does not
// care where the data comes from (userstore, CRM, external API, etc.).
//
// All methods are optional — a nil implementation means customer
// fields stay empty and customer filtering is disabled.
type CustomerResolverInterface interface {
    // FindByID returns customer display name and email for a given ID.
    // Returns empty strings if not found — no error for "not found".
    // Called by order_details and order_manager controllers.
    FindByID(ctx context.Context, customerID string) (name, email string)

    // SearchIDs returns customer IDs matching the given name and/or
    // email substrings. Empty string means "no filter on that field".
    // Called by order_manager controller.
    SearchIDs(ctx context.Context, name, email string) ([]string, error)
}
```

`AdminOptions` takes a `CustomerResolver` field:

```go
type AdminOptions struct {
    Store           shopstore.StoreInterface
    CacheStore      cachestore.StoreInterface
    Logger          *slog.Logger
    CustomerResolver CustomerResolverInterface // optional — nil is safe
    // ...
}
```

### Why an interface, not function fields

- **Grows naturally** — adding a method (e.g. `FindByEmail`, `Count`,
  `GetOrders`) is just adding to the interface. Function fields would
  require adding more fields to `AdminOptions` each time.
- **Testable** — easy to mock with a dedicated type in tests
- **Reusable** — a host project can implement it once and share across
  multiple admin instances or other modules
- **Self-documenting** — the interface is a named contract, not scattered
  function fields on a config struct
- **Nil-safe** — shopadmin nil-checks `CustomerResolver` before calling,
  so it's still optional

---

## 3. Impact

### `go.mod`

`github.com/dracory/userstore` is **removed** from `go.mod` and all imports.

### `AdminOptions`

```diff
 type AdminOptions struct {
     Store       shopstore.StoreInterface
-    UserStore   userstore.StoreInterface
     CacheStore  cachestore.StoreInterface
     Logger      *slog.Logger
+    CustomerResolver CustomerResolverInterface // optional — nil is safe
     FuncLayout      func(...) string
     AdminHomeURL    string
     ShopAdminURL    string
     FileManagerURL  string
     AuthUserID      func(r *http.Request) string
 }
```

### `RegistryInterface`

`GetUserStore()` is removed. `RegistryInterface` no longer references
`userstore`:

```diff
 type RegistryInterface interface {
     GetShopStore() shopstore.StoreInterface
-    GetUserStore() userstore.StoreInterface
     GetCacheStore() cachestore.StoreInterface
     GetLogger() *slog.Logger
 }
```

### `UiConfig` / `UiBase` / `UiInterface`

`UserStore` field replaced with `FindCustomer` / `SearchCustomerIDs` function
fields throughout:

```diff
 type UiConfig struct {
     Store      shopstore.StoreInterface
-    UserStore  userstore.StoreInterface
     CacheStore cachestore.StoreInterface
     Logger     *slog.Logger
+    CustomerResolver CustomerResolverInterface
     Layout     func(...) string
 }
```

### Controllers

`order_details` and `order_manager` replace direct `userStore` calls with
nil-checked `CustomerResolver` calls:

```go
// Before (order_details)
if order.GetCustomerID() != "" && controller.app.GetUserStore() != nil {
    customer, err := controller.app.GetUserStore().UserFindByID(ctx, order.GetCustomerID())
    if err == nil && customer != nil {
        customerName = customer.GetFirstName() + " " + customer.GetLastName()
        customerEmail = customer.GetEmail()
    }
}

// After (order_details)
if ui.CustomerResolver() != nil && order.GetCustomerID() != "" {
    customerName, customerEmail = ui.CustomerResolver().FindByID(ctx, order.GetCustomerID())
}
```

```go
// Before (order_manager)
if (reqBody.CustomerName != "" || reqBody.CustomerEmail != "") && controller.app.GetUserStore() != nil {
    users, err := controller.app.GetUserStore().UserList(ctx, userstore.NewUserQuery())
    // ... iterate ALL users to find matches
}

// After (order_manager)
if (reqBody.CustomerName != "" || reqBody.CustomerEmail != "") && ui.CustomerResolver() != nil {
    matchingCustomerIDs, err := ui.CustomerResolver().SearchIDs(ctx, reqBody.CustomerName, reqBody.CustomerEmail)
    // ... filter orders by matching IDs
}
```

---

## 4. Host Project Wiring

The host project (e.g. `mechestudio.com`) provides a `CustomerResolverInterface`
implementation at construction time:

```go
// Adapter — lives in the host project, not in shopadmin
type customerResolver struct {
    userStore userstore.StoreInterface
}

func (r *customerResolver) FindByID(ctx context.Context, id string) (string, string) {
    if r.userStore == nil || id == "" {
        return "", ""
    }
    u, err := r.userStore.UserFindByID(ctx, id)
    if err != nil || u == nil {
        return "", ""
    }
    return u.GetFirstName() + " " + u.GetLastName(), u.GetEmail()
}

func (r *customerResolver) SearchIDs(ctx context.Context, name, email string) ([]string, error) {
    if r.userStore == nil {
        return nil, nil
    }
    users, err := r.userStore.UserList(ctx, userstore.NewUserQuery())
    if err != nil {
        return nil, err
    }
    var ids []string
    for _, u := range users {
        fullName := strings.ToLower(u.GetFirstName() + " " + u.GetLastName())
        userEmail := strings.ToLower(u.GetEmail())
        if name != "" && !strings.Contains(fullName, strings.ToLower(name)) {
            continue
        }
        if email != "" && !strings.Contains(userEmail, strings.ToLower(email)) {
            continue
        }
        ids = append(ids, u.GetID())
    }
    return ids, nil
}
```

Wiring:

```go
admin, _ := shopadmin.New(shopadmin.AdminOptions{
    Store:            app.GetShopStore(),
    CacheStore:       app.GetCacheStore(),
    Logger:           app.GetLogger(),
    CustomerResolver: &customerResolver{userStore: app.GetUserStore()},
})
```

Or via `Routes()`:

```go
shopRoutes, err := shopadmin.Routes(app, shopadmin.AdminOptions{
    CustomerResolver: &customerResolver{userStore: app.GetUserStore()},
    AdminHomeURL:     links.Admin().Home(),
    ShopAdminURL:     links.Admin().Shop(),
})
```

---

## 5. Breaking Changes

1. **`AdminOptions.UserStore` field removed** — replaced by
   `CustomerResolver CustomerResolverInterface`. Consumers setting
   `UserStore` must provide a `CustomerResolverInterface` implementation
   instead. **Needs MD approval.**

2. **`RegistryInterface.GetUserStore()` removed** — `Routes()` consumers
   who relied on the registry for customer resolution must pass
   `CustomerResolver` via opts. **Needs MD approval.**

3. **`UiInterface.UserStore()` method removed** — replaced by
   `CustomerResolver()` method. Any subcontroller embedding `UiBase`
   gets the new method automatically. **Informational — internal only,
   no external consumers.**

---

## 6. Bonus: Fixes Pre-Existing Bug

The current `UserList` call in `handle_orders_load_ajax.go` loads **every
user in the database** to do substring filtering in Go. With
`CustomerResolverInterface.SearchIDs`, a host project with a real database
can implement this as a SQL `LIKE` query instead. The abstraction makes the
optimization natural rather than fighting against a concrete store API.

---

## 7. Files to Change

| File | Change |
|------|--------|
| `go.mod` | Remove `github.com/dracory/userstore` |
| `types.go` (new) | Define `CustomerResolverInterface` |
| `shopadmin.go` | Replace `UserStore` field with `CustomerResolver CustomerResolverInterface` in `AdminOptions` and `admin` struct |
| `registry.go` | Remove `GetUserStore()` from `RegistryInterface`, remove `userstore` import |
| `routes.go` | Replace `UserStore: registry.GetUserStore()` with `CustomerResolver: options.CustomerResolver` |
| `shared/ui_config.go` | Replace `UserStore` with `CustomerResolver CustomerResolverInterface` |
| `shared/ui_base.go` | Replace `UserStoreField`/`UserStore()` with `CustomerResolverField`/`CustomerResolver()` |
| `shared/ui_interface.go` | Replace `UserStore()` with `CustomerResolver()` |
| `order_details/` (when ported) | Use `CustomerResolver().FindByID()` instead of `UserStore.UserFindByID()` |
| `order_manager/` (when ported) | Use `CustomerResolver().SearchIDs()` instead of `UserStore.UserList()` |

---

## 8. Definition of Done

- [x] `CustomerResolverInterface` defined with `FindByID` and `SearchIDs` methods
- [x] `userstore` removed from `go.mod` and all imports
- [x] `CustomerResolver` wired through `AdminOptions` → `UiConfig` → `UiInterface`
- [x] `RegistryInterface` no longer references `userstore`
- [x] `go build ./...` passes
- [x] `go test ./...` passes
