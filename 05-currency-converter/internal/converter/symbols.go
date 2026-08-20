package converter

import (
	"strconv"
	"strings"
)

var CurrencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"INR": "₹",
}

func FormatAmount(val float64) string {
	s := strconv.FormatFloat(val, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}
