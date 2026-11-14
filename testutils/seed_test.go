package testutils_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

func TestSeedDB(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	testutils.SeedDB(store, slog.New(slog.NewTextHandler(os.Stderr, nil)))

	ctx := context.Background()
	if n, _ := store.ProductCount(ctx, shopstore.NewProductQuery()); n != 60 {
		t.Fatalf("expected 60 products, got %d", n)
	}
	if n, _ := store.CategoryCount(ctx, shopstore.NewCategoryQuery()); n != 40 {
		t.Fatalf("expected 40 categories, got %d", n)
	}
	if n, _ := store.DiscountCount(ctx, shopstore.NewDiscountQuery()); n != 40 {
		t.Fatalf("expected 40 discounts, got %d", n)
	}
	if n, _ := store.OrderCount(ctx, shopstore.NewOrderQuery()); n != 30 {
		t.Fatalf("expected 30 orders, got %d", n)
	}
}
