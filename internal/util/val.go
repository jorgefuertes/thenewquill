package util

import (
	"fmt"
)

const (
	FalseValue = "0"
	TrueValue  = "1"
)

func ValueToString(value any) string {
	switch val := value.(type) {
	case string:
		return val
	case bool:
		if val {
			return TrueValue
		}

		return FalseValue
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%.4f", val)
	default:
		return fmt.Sprint(value)
	}
}
