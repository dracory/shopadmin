package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dracory/shopstore"
	"github.com/dromara/carbon/v2"
)

// seedDB populates the store with sample data so the admin panel has
// enough rows to exercise pagination, filtering, and sorting.
//
// Creates:
//   - 60 products (mixed statuses, varied prices/quantities)
//   - 40 discounts (mixed statuses/types/amounts)
//   - 30 orders (mixed statuses, each with 1-3 line items)
//   - 40 categories (mixed statuses)
func seedDB(store shopstore.StoreInterface, logger *slog.Logger) {
	ctx := context.Background()

	seedCategories(ctx, store, logger)
	products := seedProducts(ctx, store, logger)
	seedDiscounts(ctx, store, logger)
	seedOrders(ctx, store, logger, products)

	logger.Info("seedDB complete",
		"products", len(products),
		"discounts", 40,
		"orders", 30,
		"categories", 40,
	)
}

func seedCategories(ctx context.Context, store shopstore.StoreInterface, logger *slog.Logger) {
	statuses := []string{"active", "active", "active", "inactive", "draft"}

	for i := 1; i <= 40; i++ {
		status := statuses[(i-1)%len(statuses)]
		cat := shopstore.NewCategory()
		cat.SetTitle(fmt.Sprintf("Category %02d", i))
		cat.SetDescription(fmt.Sprintf("Description for category %02d — used for testing pagination and filtering.", i))
		cat.SetStatus(status)

		if err := store.CategoryCreate(ctx, cat); err != nil {
			logger.Error("seedCategories: failed to create category", "index", i, "error", err)
		}
	}
}

func seedProducts(ctx context.Context, store shopstore.StoreInterface, logger *slog.Logger) []shopstore.ProductInterface {
	statuses := []string{"active", "active", "active", "inactive", "draft"}
	var products []shopstore.ProductInterface

	for i := 1; i <= 60; i++ {
		status := statuses[(i-1)%len(statuses)]
		price := fmt.Sprintf("%.2f", float64(i)*1.99)
		qty := int64(i % 50)

		p := shopstore.NewProduct()
		p.SetTitle(fmt.Sprintf("Product %02d", i))
		p.SetDescription(fmt.Sprintf("<p>Detailed description for product %02d. This is a sample product used for testing the product manager pagination, filtering, and sorting features.</p>", i))
		p.SetShortDescription(fmt.Sprintf("Short description for product %02d.", i))
		p.SetPrice(price)
		p.SetQuantity(fmt.Sprintf("%d", qty))
		p.SetStatus(status)
		p.SetMemo(fmt.Sprintf("Internal memo for product %02d", i))

		if err := store.ProductCreate(ctx, p); err != nil {
			logger.Error("seedProducts: failed to create product", "index", i, "error", err)
			continue
		}
		products = append(products, p)
	}

	return products
}

func seedDiscounts(ctx context.Context, store shopstore.StoreInterface, logger *slog.Logger) {
	statuses := []string{"active", "active", "inactive", "draft"}
	types := []string{"percent", "amount"}

	for i := 1; i <= 40; i++ {
		status := statuses[(i-1)%len(statuses)]
		discType := types[(i-1)%len(types)]

		d := shopstore.NewDiscount()
		d.SetCode(fmt.Sprintf("DISC%03d", i))
		d.SetTitle(fmt.Sprintf("Discount %02d", i))
		d.SetDescription(fmt.Sprintf("Description for discount %02d — used for testing pagination and filtering.", i))
		d.SetType(discType)

		if discType == "percent" {
			d.SetAmount(float64(i % 50))
		} else {
			d.SetAmount(float64(i) * 2.50)
		}

		d.SetStatus(status)
		d.SetStartsAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
		d.SetEndsAt(carbon.Now(carbon.UTC).AddDays(30).ToDateTimeString(carbon.UTC))
		d.SetMemo(fmt.Sprintf("Memo for discount %02d", i))

		if err := store.DiscountCreate(ctx, d); err != nil {
			logger.Error("seedDiscounts: failed to create discount", "index", i, "error", err)
		}
	}
}

func seedOrders(ctx context.Context, store shopstore.StoreInterface, logger *slog.Logger, products []shopstore.ProductInterface) {
	statuses := []string{"pending", "pending", "completed", "completed", "cancelled", "shipped"}

	for i := 1; i <= 30; i++ {
		status := statuses[(i-1)%len(statuses)]
		customerID := fmt.Sprintf("cust_%03d", (i%10)+1)

		o := shopstore.NewOrder()
		o.SetCustomerID(customerID)
		o.SetStatus(status)
		o.SetPriceFloat(float64(i) * 12.50)
		o.SetQuantityInt(int64((i % 3) + 1))
		o.SetMemo(fmt.Sprintf("Order memo %02d", i))

		// Add customer name/email as meta so order_manager can display them
		_ = o.SetMeta("customer_name", fmt.Sprintf("Customer %d", (i%10)+1))
		_ = o.SetMeta("customer_email", fmt.Sprintf("customer%d@example.com", (i%10)+1))

		if err := store.OrderCreate(ctx, o); err != nil {
			logger.Error("seedOrders: failed to create order", "index", i, "error", err)
			continue
		}

		// Add 1-3 line items per order
		itemCount := (i % 3) + 1
		for j := 0; j < itemCount; j++ {
			prodIdx := (i + j) % len(products)
			if prodIdx >= len(products) {
				prodIdx = 0
			}
			prod := products[prodIdx]

			item := shopstore.NewOrderLineItem()
			item.SetOrderID(o.GetID())
			item.SetProductID(prod.GetID())
			item.SetTitle(prod.GetTitle())
			item.SetPrice(prod.GetPrice())
			item.SetQuantityInt(int64(j + 1))
			item.SetStatus(status)

			if err := store.OrderLineItemCreate(ctx, item); err != nil {
				logger.Error("seedOrders: failed to create line item", "order", i, "item", j, "error", err)
			}
		}
	}
}
