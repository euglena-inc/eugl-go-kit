package dbvalue

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type AmountCents int64

func NewAmountCents(cents int64) AmountCents {
	return AmountCents(cents)
}

func (a AmountCents) Cents() int64 {
	return int64(a)
}

func (a AmountCents) Value() (driver.Value, error) {
	cents := int64(a)
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100), nil
}

func (a *AmountCents) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*a = 0
		return nil
	case int64:
		*a = AmountCents(v * 100)
		return nil
	case float64:
		*a = AmountCents(math.Round(v * 100))
		return nil
	case []byte:
		return a.scanString(string(v))
	case string:
		return a.scanString(v)
	default:
		return fmt.Errorf("unsupported amount type %T", value)
	}
}

func (a *AmountCents) scanString(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		*a = 0
		return nil
	}
	negative := strings.HasPrefix(value, "-")
	value = strings.TrimPrefix(value, "-")
	parts := strings.SplitN(value, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return err
	}
	var fraction int64
	if len(parts) == 2 {
		frac := parts[1]
		if len(frac) > 2 {
			frac = frac[:2]
		}
		for len(frac) < 2 {
			frac += "0"
		}
		fraction, err = strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return err
		}
	}
	cents := whole*100 + fraction
	if negative {
		cents = -cents
	}
	*a = AmountCents(cents)
	return nil
}
