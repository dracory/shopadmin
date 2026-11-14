# Proposal: Order Notifications

**Date:** 2026-11-14
**Status:** Draft — awaiting MD/GOD approval
**Modules:** `github.com/dracory/shopstore`, `github.com/dracory/shopadmin`, host applications
**Host reference:** `mechestudio.com`

---

## 1. Problem

The shop sends no application-generated order emails:

- no customer confirmation after successful payment
- no admin notification for a paid order
- no admin notification for a cancelled order
- no admin notification when checkout fails while creating order line items

Because these packages are shared by multiple projects, `shopadmin` and
`shopstore` must provide reusable event, template, and settings behavior
without owning SMTP, queues, users, branding, or application URLs.

---

## 2. Decisions

| Area | Decision |
|------|----------|
| Order events | Event types and `OrderEvent` defined in `shopstore` |
| Notification catalog | Defined in `shopstore/notifications` |
| Templates | Package defaults; per-template overrides arrive in Phase 2 |
| Admin UI | `Emails` tab inside the existing Settings page (Phase 2) |
| SMTP / tasks / users / links | Host-provided |
| Duplicate protection | `(event, audience)` sent timestamp in `order.metas` |
| Reply-To | Passed through the host `internal/emails` wrapper |
| Admin recipients | Shared `shop.notification.admin_recipients` setting |
| Line-item failure | Emits `order.cancelled`; Phase 1 notifies admin only |

---

## 3. Scope

### Phase 1 — transactional emails

- `shopstore` order event types and `OrderEvent`
- `shopstore/notifications` catalog, default templates, and renderer
- Host adapter and task that send one notification per audience
- Events for:
  - `order.paid`
  - `order.cancelled`
- Emails for:
  - customer paid-order confirmation
  - admin paid-order notification
  - admin cancelled-order notification
  - admin notification for checkout line-item failure

Admin recipients are read from the shared setting. The setting is set through the
existing settings API. No admin UI is required in Phase 1.

### Phase 2 — Emails settings screen

- `Settings → Emails` tab
- notification list and enabled toggles
- shared default admin recipients field
- per-template overrides (enabled, subject, body) stored as JSON in `Setting`
- template editor with save, restore default, preview, and a test-send modal
- malformed override fallback

### Phase 3 — admin-originated notifications

- resend receipt
- order status-change email
- any future customer cancellation or refund templates
- shared event handler interface and dispatcher in `shopstore`, added when `shopadmin`
  first emits events

Do not add event types until an emitter exists.

---

## 4. `shopstore` API

### 4.1 Order events

New `order_events.go`:

```go
package shopstore

type OrderEventType string

const (
    OrderEventPaid      OrderEventType = "order.paid"
    OrderEventCancelled OrderEventType = "order.cancelled"
)

// OrderEvent describes a persisted order lifecycle event.
type OrderEvent struct {
    Type   OrderEventType
    Order  OrderInterface
    Source string
}
```

`Source` identifies the emitter. It is used for logs and task parameters:

- `payment_success`
- `checkout` (free orders)
- `payment_cancelled`
- `checkout_line_item_failure`

No handler interface or dispatcher is defined yet. Each emitter calls the host notifier
directly. The shared interface is added to `shopstore` when `shopadmin` first emits events,
so the API is not designed ahead of its use.

### 4.2 Notifications subpackage

New package: `shopstore/notifications`.

It uses only standard-library packages plus `shopstore` types. It does not send email,
enqueue tasks, write settings, resolve users, or construct links.

| File | Purpose |
|------|---------|
| `catalog.go` | Audiences, definitions, and setting keys |
| `template.go` | Template model and validation |
| `render.go` | Default lookup, recipients, Reply-To, placeholders, rendering |
| `templates/` | Embedded default HTML |

### 4.3 Notification model

```go
package notifications

type Audience string

const (
    AudienceCustomer Audience = "customer"
    AudienceAdmin    Audience = "admin"
)

type Definition struct {
    Event      shopstore.OrderEventType
    Audience   Audience
    SettingKey string
    Title      string
    Default    Template
}

// Template is the overridable content. In Phase 2 it is stored as JSON in Setting.value.
type Template struct {
    Enabled  bool   `json:"enabled"`
    Subject  string `json:"subject"`
    BodyHTML string `json:"body_html"`
}

type Rendered struct {
    Enabled    bool
    Recipients []string
    ReplyTo    []string
    Subject    string
    BodyHTML   string
}

type SettingGetter interface {
    SettingGet(ctx context.Context, key string, defaultValue string) (string, error)
}

type RenderInput struct {
    Event         shopstore.OrderEvent
    Audience      Audience
    LineItems     []shopstore.OrderLineItemInterface
    SiteName      string
    StoreEmail    string
    Currency      string
    OrderURL      string
    AdminURL      string
    CustomerName  string
    CustomerEmail string
}
```

### 4.4 Settings

Phase 1 reads one setting:

```text
shop.notification.admin_recipients     comma-separated addresses, shared by all admin emails
```

Phase 2 adds per-template overrides, one JSON `Template` per key:

```text
shop.notification.order_paid.customer
shop.notification.order_paid.admin
shop.notification.order_cancelled.admin
```

A missing override uses the package default. All default templates ship with
`Enabled: true`. Customer cancellation emails can be added later with
`shop.notification.order_cancelled.customer`.

### 4.5 Rendering behavior

Phase 1, `RenderOrderEmail`:

1. Finds the catalog definition for `(event, audience)`.
2. Uses the default template.
3. Applies placeholders to `Subject` and `BodyHTML`.
4. Resolves recipients:
   - customer: `CustomerEmail`
   - admin: `shop.notification.admin_recipients`, parsed and validated
5. Resolves Reply-To:
   - admin: `CustomerEmail`
   - customer: `StoreEmail`
6. Returns `Enabled: false` when the template is disabled or no recipient exists.

Phase 2 adds: read the override through `SettingGetter`. If it is malformed or fails
`Validate`, log a warning and use the default.

Supported placeholders:

```text
{{.SiteName}}
{{.StoreEmail}}
{{.OrderID}}
{{.OrderNumber}}
{{.CustomerName}}
{{.CustomerEmail}}
{{.ItemsTable}}
{{.Total}}
{{.Currency}}
{{.OrderURL}}
{{.AdminURL}}
```

Safety:

- `BodyHTML` uses `html/template`. Dynamic values are escaped.
- `ItemsTable` is generated as safe HTML from escaped cell values.
- `Subject` uses `text/template`, and placeholder values have CR and LF removed to prevent
  header injection.
- `Validate` parses `Subject` and `BodyHTML`. It runs when overrides are saved in Phase 2
  and at render time.
- Admin-authored body HTML is trusted template input. Dynamic values are still escaped.
- The plain-text part is generated from HTML by `dracory/email` when `TextBody` is empty.
  No text template is maintained.

---

## 5. `shopadmin` UI (Phase 2)

### 5.1 Location

Add an **Emails** tab to Settings:

- **General** — existing generic settings manager
- **Emails** — notification management

The General tab hides `shop.notification.*` so admins do not edit raw JSON.

### 5.2 Visibility

The Emails tab appears only when `AdminOptions.NotificationSender` is set.
This keeps email UI hidden in applications without outbound mail.

### 5.3 Notification list

| Column | Source |
|--------|--------|
| Email | Definition title, linked to edit page |
| Type | `HTML` |
| Recipients | `Customer`, or the shared admin recipients |
| Enabled | AJAX toggle |

The list also contains a shared **Default admin recipients** field backed by
`shop.notification.admin_recipients`.

### 5.4 Template editor

Fields:

- enabled
- subject
- HTML body
- placeholder reference

Actions:

- Save
- Restore default
- Preview
- Send test email

### 5.5 Preview

The preview renders a built-in sample order inside:

```html
<iframe sandbox="" srcdoc="...">
```

This prevents admin-authored HTML from running scripts in the admin session.

### 5.6 Test send

**Send test email** opens a modal with one field, **Recipient**.

- the address is validated with `net/mail.ParseAddress`
- the notification is rendered with sample order data
- the subject is prefixed with `[TEST]`
- the request uses the same admin POST protection as other shopadmin actions
- nothing from the modal is persisted

`shopadmin` never sends SMTP. It calls a host-provided interface:

```go
type NotificationSenderInterface interface {
    SendRenderedEmail(
        ctx context.Context,
        to []string,
        rendered notifications.Rendered,
    ) error
}
```

Add to `AdminOptions`:

```go
NotificationSender NotificationSenderInterface
```

Because the recipient is free text, any admin can send the current template to any
address through the store's SMTP. Admins already control template content, and the `[TEST]`
prefix marks the mail. If this needs tighter control, cap test sends per admin in the
handler; that is not part of this proposal.

### 5.7 UI tests

- notification list
- enabled toggle
- save and reload
- restore default
- sandboxed preview
- test-send modal: invalid address rejected, valid address sends `[TEST]` mail

Use chromedp + `httptest`, matching the existing `shopadmin` e2e setup.

---

## 6. Host integration

### 6.1 Event notifier

New `internal/adapters/order_event_handler.go`.

For each event, it finds the catalog definitions for `event.Type` and enqueues one task per
audience. It logs enqueue errors and does not return them, so checkout and admin writes are
not affected.

Task parameters:

```go
map[string]any{
    "order_id":   orderID,
    "event_type": string(event.Type),
    "audience":   string(audience),
    "source":     event.Source,
}
```

### 6.2 Emit points

| Event | Location | When | Source |
|-------|----------|------|--------|
| `order.paid` | `payment_success_controller.go` | after paid order `OrderUpdate` succeeds | `payment_success` |
| `order.paid` | `payment_begin_controller.go` | after free order is created completed | `checkout` |
| `order.cancelled` | `payment_canceled_controller.go` | after cancelled `OrderUpdate` succeeds | `payment_cancelled` |
| `order.cancelled` | `payment_begin_controller.go` | after line-item failure cancellation | `checkout_line_item_failure` |

Phase 1 has only an admin definition for `order.cancelled`, so the customer receives no
email for the internal line-item failure.

### 6.3 Email task

New `internal/tasks/order_email_delivery/`.

Each task handles exactly one `(order, event, audience)` pair:

1. Load order, customer, line items, and URLs.
2. Check `notification.sent.<event>.<audience>` in `order.metas`.
3. Skip when already sent.
4. Render through `notifications.RenderOrderEmail`.
5. Skip and log when disabled or no recipient exists.
6. Send through `emails.SendEmail`.
7. Write the sent timestamp to `order.metas`.

Example send call:

```go
err := emails.SendEmail(emails.SendOptions{
    From:     app.GetConfig().GetMailFromAddress(),
    FromName: app.GetConfig().GetMailFromName(),
    To:       rendered.Recipients,
    ReplyTo:  rendered.ReplyTo,
    Subject:  rendered.Subject,
    HtmlBody: rendered.BodyHTML,
})
```

`internal/emails.SendOptions` supports `ReplyTo []string` and passes it to `dracory/email`,
which already accepts it.

### 6.4 Duplicate protection

Successful sends write metadata such as:

```json
{
  "notification.sent.order_paid.customer": "2026-11-14T12:34:56Z",
  "notification.sent.order_paid.admin": "2026-11-14T12:34:58Z"
}
```

A retry checks this marker before sending.

This is intentionally simple. A repeated Stripe callback may enqueue a duplicate task,
but the second task exits after the first marks the audience as sent. Metadata updates
are not atomic, so a rare concurrent task or failed metadata write could still duplicate
an email. We accept that risk for Phase 1 and can revisit it only if observed.

### 6.5 Admin recipients

Seed `shop.notification.admin_recipients` through the existing once-seeder pattern in
`database/seeders/once/settings`. The recipient list is supplied per environment, not
hard-coded in application code. Without it, admin notifications are skipped and logged.

---

## 7. Failure behavior

- Enqueue errors are logged and do not fail order updates.
- Missing admin recipients: admin emails are skipped and logged.
- Missing customer email: the customer email is skipped and logged.
- SMTP failures fail the task and use task-store retry behavior.
- Retries check the metadata marker before attempting delivery.
- Phase 2: a malformed override falls back to the default and logs a warning.

---

## 8. Compatibility

- New `shopstore` types and subpackage; no existing signatures change.
- New optional `AdminOptions` fields; nil preserves current behavior.
- Uses the existing `Setting` entity; no new table.
- Hosts without `NotificationSender` see no Emails UI.
- Projects with databases older than `Setting` need the normal setting-table migration
  before using this feature.

---

## 9. Files

### `shopstore`

| File | Change |
|------|--------|
| `order_events.go` | Event types and `OrderEvent` |
| `notifications/catalog.go` | Audiences, definitions, setting keys |
| `notifications/template.go` | Template model and validation |
| `notifications/render.go` | Recipients, Reply-To, placeholders, rendering |
| `notifications/templates/` | Default HTML templates |

### `shopadmin` (Phase 2)

| File | Change |
|------|--------|
| `shopadmin.go` | `NotificationSender` option |
| `setting_manager/` | Settings tabs; hide notification keys from General |
| `email_notifications/` | List, editor, preview, and test-send UI |
| `example/main.go` | Example sender wiring |

### `mechestudio.com`

| File | Change |
|------|--------|
| `internal/adapters/order_event_handler.go` | Enqueue notification tasks |
| `internal/adapters/notification_sender.go` | Test-send sender (Phase 2) |
| `internal/tasks/order_email_delivery/` | Render, send, mark sent |
| `internal/tasks/constants/constants.go` | Add task alias |
| `internal/tasks/register_tasks.go` | Register task |
| `internal/controllers/website/payment/payment_success_controller.go` | Emit `order.paid` |
| `internal/controllers/website/payment/payment_canceled_controller.go` | Emit `order.cancelled` |
| `internal/controllers/user/order/payment_begin_controller.go` | Emit free-order paid and line-item failure cancelled events |
| `database/seeders/once/settings/` | Seed `shop.notification.admin_recipients` |

`internal/emails/send_email.go` already supports `ReplyTo`.

---

## 10. Definition of Done

### Phase 1

- [ ] `shopstore` event types and `OrderEvent` implemented
- [ ] notification catalog and renderer implemented
- [ ] default templates render
- [ ] customer paid confirmation sends
- [ ] admin paid notification sends
- [ ] admin cancelled notification sends
- [ ] line-item failure sends admin-only cancelled notification
- [ ] free orders emit `order.paid`
- [ ] admin recipients read from setting; missing setting skips and logs
- [ ] Reply-To set correctly for both audiences
- [ ] successful sends write metadata timestamps
- [ ] retries skip already-sent `(event, audience)` deliveries
- [ ] enqueue failures do not break checkout
- [ ] tests cover events, rendering, escaping, recipients, Reply-To, and task behavior

### Phase 2

- [ ] Emails tab appears only when sender is configured
- [ ] shared admin recipients can be edited
- [ ] list, toggle, edit, save, restore, and preview work
- [ ] overrides apply; malformed overrides fall back to defaults with a warning
- [ ] preview is sandboxed
- [ ] test-send modal validates the address and sends a `[TEST]` email
- [ ] General tab hides notification keys
- [ ] e2e covers toggle and save

---

## 11. Release order

1. Release `shopstore`.
2. Release `shopadmin` against the new `shopstore` version.
3. Bump `mechestudio.com` dependencies.
4. Ensure local `replace` directives remain commented before `task deploy-live`.
