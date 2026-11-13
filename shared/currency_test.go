package shared

import "testing"

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
