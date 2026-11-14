package shared

import (
	"context"
	"testing"

	"github.com/dracory/shopadmin/testutils"

	_ "modernc.org/sqlite"
)

func TestFormatPrice(t *testing.T) {
	if got := FormatPrice("12.5", "GBP"); got != "£12.50" {
		t.Fatalf("expected £12.50, got %q", got)
	}
	if got := FormatPrice(12, "USD"); got != "$12.00" {
		t.Fatalf("expected $12.00, got %q", got)
	}
	if got := FormatPrice("9.999", "EUR"); got != "€9.99" {
		t.Fatalf("expected €9.99, got %q", got)
	}
	if got := FormatPrice("7.5", "XYZ"); got != "XYZ 7.50" {
		t.Fatalf("expected XYZ 7.50, got %q", got)
	}
}

func TestCurrencySymbol(t *testing.T) {
	if CurrencySymbol("gbp") != "£" {
		t.Fatalf("expected £, got %q", CurrencySymbol("gbp"))
	}
	if CurrencySymbol("XXX") != "XXX " {
		t.Fatalf("expected 'XXX ', got %q", CurrencySymbol("XXX"))
	}
}

func TestShopCurrencyNilStore(t *testing.T) {
	if got := ShopCurrency(t.Context(), nil); got != CURRENCY_DEFAULT {
		t.Fatalf("expected %q, got %q", CURRENCY_DEFAULT, got)
	}
}

func TestShopCurrency(t *testing.T) {
	if got := ShopCurrency(context.Background(), nil); got != CURRENCY_DEFAULT {
		t.Fatalf("expected default currency for nil store, got %q", got)
	}

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("init store: %v", err)
	}
	ctx := context.Background()

	if got := ShopCurrency(ctx, store); got != CURRENCY_DEFAULT {
		t.Fatalf("expected default currency without setting, got %q", got)
	}

	if err := store.SettingSet(ctx, SETTING_KEY_CURRENCY, "gbp"); err != nil {
		t.Fatalf("set currency: %v", err)
	}
	if got := ShopCurrency(ctx, store); got != "GBP" {
		t.Fatalf("expected GBP, got %q", got)
	}
}
