package product_update_test

import (
	"context"

	"github.com/dracory/shopstore"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// postForm posts a form-encoded body with no action param, which
// drives the component Mount/Handle/Render path.
func postForm(t *testing.T, env *testEnv, query url.Values, form url.Values) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=product_update&"+query.Encode(),
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)
	return rec
}

func getView(t *testing.T, env *testEnv, productID, view string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet,
		"/admin/shop?controller=product_update&product_id="+productID+"&view="+view,
		http.NoBody)
	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)
	return rec
}

func TestPostViewDetailsComponent(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := postForm(t, env,
		url.Values{"product_id": {p.GetID()}, "view": {"details"}},
		url.Values{})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestPostViewMetadataComponent(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := postForm(t, env,
		url.Values{"product_id": {p.GetID()}, "view": {"metadata"}},
		url.Values{})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestPostViewTagsComponent(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := postForm(t, env,
		url.Values{"product_id": {p.GetID()}, "view": {"tags"}},
		url.Values{})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestPostViewMediaComponent(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := postForm(t, env,
		url.Values{"product_id": {p.GetID()}, "view": {"media"}},
		url.Values{
			"media_url":  {"https://example.com/a.png"},
			"media_type": {"image/png"},
		})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

}

func TestPostViewMediaCountMismatch(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := postForm(t, env,
		url.Values{"product_id": {p.GetID()}, "view": {"media"}},
		url.Values{
			"media_url": {"https://example.com/a.png", "https://example.com/b.png"},
		})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetViewMedia(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := getView(t, env, p.GetID(), "media")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "media") {
		t.Fatal("expected media view rendered")
	}
}

func TestGetViewTags(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := getView(t, env, p.GetID(), "tags")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetViewMetadata(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := getView(t, env, p.GetID(), "metadata")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetViewMediaWithExistingMedia(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	media := shopstore.NewMedia()
	media.SetEntityID(p.GetID())
	media.SetURL("https://example.com/pic.png")
	media.SetTitle("pic.png")
	media.SetType("image/png")
	media.SetStatus(shopstore.MEDIA_STATUS_ACTIVE)
	media.SetSequence(1)
	if err := env.store.MediaCreate(context.Background(), media); err != nil {
		t.Fatalf("seed media: %v", err)
	}

	rec := getView(t, env, p.GetID(), "media")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
