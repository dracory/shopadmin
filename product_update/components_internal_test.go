package product_update

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/shopadmin/shared"
	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

// Components whose Handle is never invoked through the controller's
// POST view dispatch (status, categories, admin notes, variants are
// rendered as cards, not as standalone views) are exercised directly
// here to cover their Handle/Mount/Render surface.

func newInternalUI(t *testing.T) (*ui, shopstore.ProductInterface) {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	u := UI(shared.UiConfig{
		Store: store,
	}, "").(*ui)

	product := shopstore.NewProduct()
	product.SetTitle("Internal")
	product.SetStatus("active")
	if err := store.ProductCreate(context.Background(), product); err != nil {
		t.Fatalf("seed product: %v", err)
	}

	return u, product
}

func TestStatusComponentHandle(t *testing.T) {
	u, product := newInternalUI(t)
	r := httptest.NewRequest("POST", "/", nil)

	c := NewProductStatusComponent(u)
	c.Mount(r, product, product.GetID())
	if err := c.Handle(r); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if c.Render() == nil {
		t.Fatal("expected non-nil render")
	}
}

func TestCategoriesComponentHandle(t *testing.T) {
	u, product := newInternalUI(t)
	r := httptest.NewRequest("POST", "/", nil)

	c := NewProductCategoriesComponent(u)
	c.Mount(r, product, product.GetID())
	if err := c.Handle(r); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if c.Render() == nil {
		t.Fatal("expected non-nil render")
	}
}

func TestAdminNotesComponentHandle(t *testing.T) {
	u, product := newInternalUI(t)
	r := httptest.NewRequest("POST", "/", nil)

	c := NewProductAdminNotesComponent(u)
	c.Mount(r, product, product.GetID())
	if err := c.Handle(r); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if c.Render() == nil {
		t.Fatal("expected non-nil render")
	}
}

func TestVariantsComponentHandle(t *testing.T) {
	u, product := newInternalUI(t)
	r := httptest.NewRequest("POST", "/", nil)

	c := NewProductVariantsComponent(u)
	c.Mount(r, product, product.GetID())
	if err := c.Handle(r); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if c.Render() == nil {
		t.Fatal("expected non-nil render")
	}
}

func TestMediaComponentCreateMediaItem(t *testing.T) {
	u, product := newInternalUI(t)
	r := httptest.NewRequest("GET", "/", nil)

	media := shopstore.NewMedia()
	media.SetEntityID(product.GetID())
	media.SetURL("https://example.com/pic.png")
	media.SetTitle("pic.png")
	media.SetType(shopstore.MEDIA_TYPE_IMAGE_JPG)

	c := NewProductMediaComponent(u)
	c.Mount(r, product, product.GetID())

	tag := c.createMediaItem(media, 0)
	if tag == nil {
		t.Fatal("expected non-nil tag")
	}
	html := tag.ToHTML()
	for _, want := range []string{"pic.png", "https://example.com/pic.png", "media_url[0]", "media_type[0]"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in media item HTML", want)
		}
	}

	// Video type renders a different icon/badge.
	video := shopstore.NewMedia()
	video.SetType(shopstore.MEDIA_TYPE_VIDEO_MP4)
	video.SetURL("https://example.com/v.mp4")
	video.SetTitle("")
	html = c.createMediaItem(video, 1).ToHTML()
	if !strings.Contains(html, "Video") || !strings.Contains(html, "Media 2") {
		t.Fatalf("expected video badge and default title, got %q", html)
	}
}
