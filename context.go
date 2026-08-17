package shopadmin

import (
	"context"

	"github.com/dracory/shopadmin/shared"
)

// contextWithValues injects shopadmin config into the request context,
// following the cmsstore/admin pattern.
func contextWithValues(ctx context.Context, opts AdminOptions, endpoint string) context.Context {
	ctx = context.WithValue(ctx, shared.KeyEndpoint, endpoint)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, opts.AdminHomeURL)
	ctx = context.WithValue(ctx, shared.KeyShopAdminURL, opts.ShopAdminURL)
	ctx = context.WithValue(ctx, shared.KeyFileManagerURL, opts.FileManagerURL)
	return ctx
}
