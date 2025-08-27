package nullify

import "github.com/shopspring/decimal"

// Empty converts empty or nil values to nil for nullable fields.
func Empty(value *string) interface{} {
	if value == nil || *value == "" {
		return nil
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
