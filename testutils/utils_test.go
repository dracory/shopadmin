package testutils_test

import (
	"context"
	"testing"

	"github.com/dracory/shopadmin/testutils"
	"github.com/dracory/shopstore"

	_ "modernc.org/sqlite"
)

func TestInitStore(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}

	// Tables should be automigrated and usable.
	p := shopstore.NewProduct()
	p.SetTitle("P")
	if err := store.ProductCreate(context.Background(), p); err != nil {
		t.Fatalf("product create: %v", err)
	}
}
