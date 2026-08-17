package testutils

import (
	"database/sql"
	"os"
	"strings"

	"github.com/dracory/shopstore"
)

// InitStore creates an in-memory shopstore for testing.
//
// Note: The SQLite driver (modernc.org/sqlite) must be imported with a
// side-effect import in _test.go files of the consuming package, NOT here.
// This avoids double driver registration panics when a host project that
// already imports modernc.org/sqlite also imports this testutils package.
//
// Example (in your _test.go file):
//
//	import (
//	    "github.com/dracory/shopadmin/testutils"
//	    _ "modernc.org/sqlite"
//	)
func InitStore(filepath string) (shopstore.StoreInterface, error) {
	db, err := initDB(filepath)
	if err != nil {
		return nil, err
	}

	store, err := shopstore.NewStore(shopstore.NewStoreOptions{
		DB:                     db,
		CategoryTableName:      "shop_category",
		DiscountTableName:      "shop_discount",
		MediaTableName:         "shop_media",
		OrderTableName:         "shop_order",
		OrderLineItemTableName: "shop_order_line_item",
		ProductTableName:       "shop_product",
		AutomigrateEnabled:     true,
	})
	if err != nil {
		return nil, err
	}

	return store, nil
}

func initDB(filepath string) (*sql.DB, error) {
	if filepath != ":memory:" {
		err := os.Remove(filepath)
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			return nil, err
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
