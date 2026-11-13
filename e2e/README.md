# e2e

End-to-end tests that drive the shopadmin panel in a real headless Chrome
via [chromedp](https://github.com/chromedp/chromedp) (CDP). Unlike go-rod,
chromedp has no `leakless.exe` helper, so it avoids the Windows Defender
false positive entirely.

## How it works

- `TestMain` mounts the admin on `httptest.NewServer` (in-process, random
  port) backed by a seeded in-memory SQLite store — no binary build needed.
- One shared browser is launched via `chromedp.NewExecAllocator` for the
  whole run; each test opens a fresh tab with a 30s deadline (`newTab`).
- Chrome binary resolution order: `CHROME_BIN` env → Windows Playwright
  `chromium_headless_shell` → WSL-side Windows shell → Linux Playwright
  shell → chromedp's default lookup (Chrome on PATH / install dirs).

## Run

```bash
go test -v ./e2e/     # or: task e2e
```

Note: on WSL, run tests with the Windows Go toolchain (`go.exe`) —
Windows-side Chrome binds 127.0.0.1, which WSL cannot reach (no
WSL→Windows localhost forwarding).

## Gotchas learned

- Prefer `chromedp.WaitReady` over `chromedp.WaitVisible`: on some pages
  `WaitVisible` polls forever even though the DOM is fully loaded.
- Routing is query-param based: `?controller=products|categories|...`
  (see `shared/consts.go`).
- Manager tables render via Vue after an AJAX fetch — the raw `v-for`
  template exists in the DOM before mount, so wait for
  `table tbody tr:nth-of-type(2)` (a second row only exists post-render).
- Vue strips directives like `v-model` and `:id` at mount — they can't
  be used as selectors. Use `#app[data-v-app]`, `[id^="cat-"]`, or
  `input[disabled]` as post-mount markers, and `chromedp.Poll` to wait
  for prefilled values.
- `chromedp.Poll` waits for a JS predicate — use it for "value loaded"
  conditions where `WaitReady` only proves the element exists.
- `customer_name` resolves via `CustomerResolver`, which is nil here —
  seeded orders show `-`; assert on status instead.
- chromedp v0.23 rewrote the API (generic `Run[T]`); this suite pins the
  stable v0.14 line with matching `cdproto`/`go-json-experiment` versions.
