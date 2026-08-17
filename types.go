package shopadmin

import "github.com/dracory/shopadmin/shared"

// CustomerResolverInterface resolves customer data for order views.
// The host project provides an implementation — shopadmin does not
// care where the data comes from (userstore, CRM, external API, etc.).
//
// This is a re-export of shared.CustomerResolverInterface so consumers
// of the root shopadmin package can use shopadmin.CustomerResolverInterface
// without importing the shared subpackage.
type CustomerResolverInterface = shared.CustomerResolverInterface
