package nullify

import "github.com/shopspring/decimal"

// EmptyString converts empty or nil values to nil for nullable fields.
func EmptyString(value *string) interface{} {
	if value == nil || *value == "" {
		return nil
	}

	return *value
}

// EmptyInt64 returns nil if the pointer is nil; otherwise it returns the
// dereferenced int64 value. Useful when passing values to layers (e.g., DB)
// that interpret nil as NULL.
func EmptyInt64(value *int64) interface{} {
	if value == nil {
		return nil
	}

	return *value
}

// EmptyInt64WithDefault returns a default int64 value when the pointer is nil;
// otherwise it returns the dereferenced int64 value. This is handy when you
// always want a concrete number, even if the input is absent.
func EmptyInt64WithDefault(value *int64) interface{} {
	const defValue int64 = 0

	if value == nil {
		return defValue
	}

	return *value
}

// EmptyDecimal handles nullable decimal values and converts them to nil if empty or nil.
func EmptyDecimal(value *decimal.Decimal) interface{} {
	if value == nil {
		return nil
	}

	return *value
}

// EmptyDecimalWithDefault handles nullable decimal values and returns a default value if empty or nil.
func EmptyDecimalWithDefault(value *decimal.Decimal) interface{} {
	const defValue = 0.00

	if value == nil {
		return decimal.NewFromFloat(defValue)
	}

	return *value
}
