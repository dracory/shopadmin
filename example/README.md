# shopadmin example

A minimal, runnable server that mounts the shopadmin panel on an
in-memory SQLite database. No external services required.

## Run

From the `shopadmin` module root:

```bash
go run ./example
```

Then open <http://localhost:8080/> in your browser and click
**Open Shop Admin**, or go directly to <http://localhost:8080/admin/shop>.

## What you get

- `/` — landing page with a link into the admin
- `/admin/shop` — shop admin dashboard
- `/admin/shop?controller=products` — product manager
- `/admin/shop?controller=categories` — category manager
- `/admin/shop?controller=discounts` — discount manager
- `/admin/shop?controller=orders` — order manager

## Configuration

Edit the constants at the top of [`main.go`](main.go) to change:

- `addr` — listen address (default `:8080`)
- `dbFile` — `:memory:` for an ephemeral DB, or a file path to persist
  data across restarts
- `adminURL`, `homeURL`, `filesURL` — mount points

## Notes

- `AuthUserID` is intentionally nil, so the panel is open. Provide a
  real implementation in production.
- The in-memory database is reset on every restart. Pass a file path
  to `openDB` to persist data.
