package shopadmin

import (
	"log/slog"

	"github.com/dracory/cachestore"
	"github.com/dracory/shopstore"
)

// RegistryInterface provides access to the stores and services that
// shopadmin needs. It is structurally compatible with
// project/internal/app.AppInterface — any type implementing
// AppInterface already satisfies this narrower interface.
//
// This interface exists so that Routes() can preserve its original
// signature Routes(registry, opts...) without importing
// project/internal/app.
//
// GetUserStore() is intentionally absent — customer resolution is
// handled via FindCustomer / SearchCustomerIDs function fields in
// AdminOptions, keeping shopadmin decoupled from userstore.
type RegistryInterface interface {
	GetShopStore() shopstore.StoreInterface
	GetCacheStore() cachestore.StoreInterface
	GetLogger() *slog.Logger
}
