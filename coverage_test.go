package shopadmin

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"
	_ "modernc.org/sqlite"
)

type covEnv struct {
	admin AdminInterface
	store shopstore.StoreInterface
}

func newCovEnv(t *testing.T) *covEnv {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	return &covEnv{admin: a, store: store}
}

func covSeedProduct(t *testing.T, env *covEnv) shopstore.ProductInterface {
	t.Helper()
	p := shopstore.NewProduct()
	p.SetTitle("Cov Product")
	p.SetStatus("active")
	p.SetPrice("9.99")
	p.SetQuantity("3")
	if err := env.store.ProductCreate(context.Background(), p); err != nil {
		t.Fatalf("seed product: %v", err)
	}
	return p
}

func covSeedCategory(t *testing.T, env *covEnv) shopstore.CategoryInterface {
	t.Helper()
	c := shopstore.NewCategory()
	c.SetTitle("Cov Category")
	c.SetStatus("active")
	if err := env.store.CategoryCreate(context.Background(), c); err != nil {
		t.Fatalf("seed category: %v", err)
	}
	return c
}

func covSeedDiscount(t *testing.T, env *covEnv) shopstore.DiscountInterface {
	t.Helper()
	d := shopstore.NewDiscount()
	d.SetCode("COV10")
	d.SetTitle("Cov Discount")
	d.SetType("percent")
	d.SetAmount(10)
	d.SetStatus("active")
	if err := env.store.DiscountCreate(context.Background(), d); err != nil {
		t.Fatalf("seed discount: %v", err)
	}
	return d
}

func covSeedOrder(t *testing.T, env *covEnv) shopstore.OrderInterface {
	t.Helper()
	o := shopstore.NewOrder()
	o.SetCustomerID("cust_1")
	o.SetStatus("pending")
	o.SetPriceFloat(19.99)
	o.SetQuantityInt(1)
	if err := env.store.OrderCreate(context.Background(), o); err != nil {
		t.Fatalf("seed order: %v", err)
	}
	return o
}

// get issues a GET through the full dispatch and asserts a 200 with content.
func covGet(t *testing.T, env *covEnv, target string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, http.NoBody)
	rr := httptest.NewRecorder()
	env.admin.Handle(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET %s: expected 200, got %d", target, rr.Code)
	}
	if rr.Body.Len() == 0 {
		t.Fatalf("GET %s: empty body", target)
	}
	return rr.Body.String()
}

type covResp struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

// covPost issues a JSON POST through the full dispatch.
func covPost(t *testing.T, env *covEnv, controller string, action string, extra url.Values, payload any) covResp {
	t.Helper()
	return covPostResp[covResp](t, env, controller, action, extra, payload)
}

// covPostResp issues a JSON POST and unmarshals the response into R.
func covPostResp[R any](t *testing.T, env *covEnv, controller string, action string, extra url.Values, payload any) R {
	t.Helper()

	q := url.Values{"controller": {controller}, "action": {action}}
	for k, vs := range extra {
		q[k] = vs
	}
	target := "/admin/shop?" + q.Encode()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	env.admin.Handle(rr, req)

	var resp R
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("POST %s: unmarshal: %v (body=%q)", target, err, rr.Body.String())
	}
	return resp
}

func covPostOK(t *testing.T, env *covEnv, controller string, action string, extra url.Values, payload any) covResp {
	t.Helper()
	resp := covPost(t, env, controller, action, extra, payload)
	if resp.Status != "success" {
		t.Fatalf("POST action=%s: expected success, got %q (message=%q)", action, resp.Status, resp.Message)
	}
	return resp
}

func TestCoverageAllPages(t *testing.T) {
	env := newCovEnv(t)

	p := covSeedProduct(t, env)
	c := covSeedCategory(t, env)
	d := covSeedDiscount(t, env)
	o := covSeedOrder(t, env)

	pages := []string{
		"/admin/shop",
		"/admin/shop?controller=products",
		"/admin/shop?controller=product_view&product_id=" + p.GetID(),
		"/admin/shop?controller=product_update&product_id=" + p.GetID(),
		"/admin/shop?controller=categories",
		"/admin/shop?controller=category_create",
		"/admin/shop?controller=category_update&category_id=" + c.GetID(),
		"/admin/shop?controller=discounts",
		"/admin/shop?controller=discount_view&discount_id=" + d.GetID(),
		"/admin/shop?controller=discount_update&discount_id=" + d.GetID(),
		"/admin/shop?controller=orders",
		"/admin/shop?controller=order_details&order_id=" + o.GetID(),
		"/admin/shop?controller=settings",
	}
	for _, target := range pages {
		covGet(t, env, target)
	}
}

func TestCoverageListAjax(t *testing.T) {
	env := newCovEnv(t)
	covSeedProduct(t, env)
	covSeedCategory(t, env)
	covSeedDiscount(t, env)
	covSeedOrder(t, env)

	actions := []struct {
		controller string
		action     string
		extra      url.Values
	}{
		{"home", "load-stats", nil},
		{"products", "load-products", nil},
		{"categories", "load-categories", nil},
		{"discounts", "load-discounts", nil},
		{"orders", "load-orders-ajax", nil},
		{"settings", "load-settings", nil},
	}
	for _, tc := range actions {
		covPostOK(t, env, tc.controller, tc.action, tc.extra, map[string]any{})
	}
}

func TestCoverageProductManagerAjax(t *testing.T) {
	env := newCovEnv(t)
	p := covSeedProduct(t, env)

	covPostOK(t, env, "products", "create-product-ajax", nil, map[string]any{"title": "New One"})
	covPostOK(t, env, "products", "delete-selected", nil, map[string]any{
		"bulk_product_ids": []string{p.GetID()},
	})
	covPostOK(t, env, "products", "delete-product", nil, map[string]any{"product_id": "nonexistent"})

	resp := covPost(t, env, "products", "load-products", nil, map[string]any{"status": "active"})
	if resp.Status != "success" {
		t.Fatalf("load-products filtered: %q", resp.Status)
	}
}

func TestCoverageCategoryAjax(t *testing.T) {
	env := newCovEnv(t)
	c := covSeedCategory(t, env)

	covPostOK(t, env, "category_create", "create-category", nil, map[string]any{"title": "Cat Two"})

	resp := covPost(t, env, "category_update", "load-category", url.Values{"category_id": {c.GetID()}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("load-category: %q (message=%q)", resp.Status, resp.Message)
	}

	covPostOK(t, env, "category_update", "update-category", nil, map[string]any{
		"category_id": c.GetID(),
		"title":       "Renamed",
		"status":      "active",
	})
	covPostOK(t, env, "categories", "delete-category", nil, map[string]any{"category_id": c.GetID()})

	c2 := covSeedCategory(t, env)
	covPostOK(t, env, "categories", "delete-selected", nil, map[string]any{
		"bulk_category_ids": []string{c2.GetID()},
	})
}

func TestCoverageDiscountAjax(t *testing.T) {
	env := newCovEnv(t)
	d := covSeedDiscount(t, env)

	covPostOK(t, env, "discounts", "create-discount-ajax", nil, map[string]any{"title": "D2", "type": "percent"})

	resp := covPost(t, env, "discount_view", "load-details", url.Values{"discount_id": {d.GetID()}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("discount load-details: %q (message=%q)", resp.Status, resp.Message)
	}

	covPostOK(t, env, "discount_update", "save-details", url.Values{"discount_id": {d.GetID()}}, map[string]any{
		"code": "COV10", "title": "Updated", "type": "percent", "amount": 15, "status": "active",
	})
	covPostOK(t, env, "discounts", "delete-discount", nil, map[string]any{"discount_id": d.GetID()})

	d2 := covSeedDiscount(t, env)
	covPostOK(t, env, "discounts", "delete-selected", nil, map[string]any{
		"bulk_discount_ids": []string{d2.GetID()},
	})
}

func TestCoverageOrderAndProductViewAjax(t *testing.T) {
	env := newCovEnv(t)
	o := covSeedOrder(t, env)
	p := covSeedProduct(t, env)

	resp := covPost(t, env, "order_details", "load-order-details", url.Values{"order_id": {o.GetID()}}, map[string]any{"order_id": o.GetID()})
	if resp.Status != "success" {
		t.Fatalf("load-order-details: %q (message=%q)", resp.Status, resp.Message)
	}

	resp = covPost(t, env, "product_view", "load-details", url.Values{"product_id": {p.GetID()}}, map[string]any{})
	if resp.Status != "success" {
		t.Fatalf("product load-details: %q (message=%q)", resp.Status, resp.Message)
	}
}

func TestCoverageProductUpdateAjax(t *testing.T) {
	env := newCovEnv(t)
	p := covSeedProduct(t, env)
	extra := url.Values{"product_id": {p.GetID()}}

	for _, action := range []string{"load-details", "load-media", "load-metadata", "load-tags"} {
		covPostOK(t, env, "product_update", action, extra, map[string]any{})
	}

	covPostOK(t, env, "product_update", "save-details", extra, map[string]any{
		"status": "active", "title": "Upd", "price": "5.00", "quantity": "2",
	})
	covPostOK(t, env, "product_update", "save-tags", extra, map[string]any{"tags": []string{"a", "b"}})
	covPostOK(t, env, "product_update", "save-metadata", extra, map[string]any{
		"metadata": []map[string]any{{"key": "k1", "value": "v1"}},
	})
	covPostOK(t, env, "product_update", "save-media", extra, map[string]any{"media": []map[string]any{}})
}

func TestCoverageProductUpdateCategoriesAjax(t *testing.T) {
	env := newCovEnv(t)
	p := covSeedProduct(t, env)
	c1 := covSeedCategory(t, env)
	c2 := covSeedCategory(t, env)
	extra := url.Values{"product_id": {p.GetID()}}

	covPostOK(t, env, "product_update", "save-categories", extra, map[string]any{
		"category_ids": []string{c1.GetID(), c2.GetID()},
	})

	resp := covPostResp[map[string]any](t, env, "product_update", "load-categories", extra, map[string]any{})
	if resp["status"] != "success" {
		t.Fatalf("load-categories: %v", resp)
	}
	selected, ok := resp["selected"].([]any)
	if !ok || len(selected) != 2 {
		t.Fatalf("expected 2 selected categories, got %v", resp["selected"])
	}

	covPostOK(t, env, "product_update", "save-categories", extra, map[string]any{
		"category_ids": []string{},
	})

	resp = covPostResp[map[string]any](t, env, "product_update", "load-categories", extra, map[string]any{})
	selected, ok = resp["selected"].([]any)
	if !ok || len(selected) != 0 {
		t.Fatalf("expected 0 selected categories after clear, got %v", resp["selected"])
	}
}

func TestCoverageSettingsAjaxViaDispatch(t *testing.T) {
	env := newCovEnv(t)

	covPostOK(t, env, "settings", "save-setting", nil, map[string]any{"key": "currency", "value": "GBP"})
	resp := covPostOK(t, env, "settings", "load-settings", nil, map[string]any{})
	if resp.Data["currency"] != "GBP" {
		t.Fatalf("expected currency=GBP, got %v", resp.Data["currency"])
	}
}
