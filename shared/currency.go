package shared

import (
	"context"
	"strings"

	"github.com/dracory/shopstore"
	"github.com/spf13/cast"
)

// SETTING_KEY_CURRENCY is the setting key used to store the shop's
// display currency (ISO 4217 code, e.g. "GBP").
const SETTING_KEY_CURRENCY = "currency"

// CURRENCY_DEFAULT is the fallback currency used when no currency
// setting has been saved or the lookup fails.
const CURRENCY_DEFAULT = "USD"

// currencySymbols maps common ISO 4217 codes to their display symbols.
var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"CNY": "¥",
	"AUD": "A$",
	"CAD": "C$",
	"NZD": "NZ$",
	"CHF": "CHF ",
	"INR": "₹",
	"BRL": "R$",
	"ZAR": "R",
	"SEK": "kr ",
	"PLN": "zł",
	"BGN": "лв",
	"RON": "lei ",
}

// ShopCurrency returns the configured shop currency code, or
// CURRENCY_DEFAULT when unset or on lookup failure.
func ShopCurrency(ctx context.Context, store shopstore.StoreInterface) string {
	if store == nil {
		return CURRENCY_DEFAULT
	}

	currency, err := store.SettingGet(ctx, SETTING_KEY_CURRENCY, CURRENCY_DEFAULT)
	if err != nil || strings.TrimSpace(currency) == "" {
		return CURRENCY_DEFAULT
	}

	return strings.ToUpper(strings.TrimSpace(currency))
}

// CurrencySymbol returns the display symbol for a currency code.
// Falls back to the code itself followed by a space (e.g. "XYZ ").
func CurrencySymbol(code string) string {
	if symbol, ok := currencySymbols[strings.ToUpper(code)]; ok {
		return symbol
	}

	return strings.ToUpper(code) + " "
}

// FormatPrice formats a numeric amount with the currency's symbol
// (e.g. "£12.00"). Unknown codes fall back to "12.00 XYZ" style.
func FormatPrice(amount any, currency string) string {
	value := cast.ToFloat64(amount)
	symbol := CurrencySymbol(currency)
	formatted := cast.ToString(value)

	parts := strings.SplitN(formatted, ".", 2)
	intPart := parts[0]
	fracPart := "00"
	if len(parts) == 2 {
		fracPart = parts[1]
		for len(fracPart) < 2 {
			fracPart += "0"
		}
		if len(fracPart) > 2 {
			fracPart = fracPart[:2]
		}
	}

	return symbol + intPart + "." + fracPart
}
