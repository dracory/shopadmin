package shopadmin

import "github.com/dracory/shopadmin/shared"

// CustomerResolverInterface resolves customer data for order views.
// The host project provides an implementation — shopadmin does not
// care where the data comes from (userstore, CRM, external API, etc.).
//
// These are re-exports from shared package so consumers of the root
// shopadmin package can use them without importing the shared subpackage.
type CustomerResolverInterface = shared.CustomerResolverInterface

type ProductSectionPosition = shared.ProductSectionPosition

const (
	PositionTop        = shared.PositionTop
	PositionBeforeForm = shared.PositionBeforeForm
	PositionAfterForm  = shared.PositionAfterForm
	PositionSidebar    = shared.PositionSidebar
	PositionBottom     = shared.PositionBottom
)

type ProductSection = shared.ProductSection
type ProductSectionRenderer = shared.ProductSectionRenderer
