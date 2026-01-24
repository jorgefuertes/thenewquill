package util

import "fmt"

func Compare(v1, v2 any) bool {
	switch v1.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return fmt.Sprintf("%d", v1) == fmt.Sprintf("%d", v2)
	case float32, float64:
		return fmt.Sprintf("%.4f", v1) == fmt.Sprintf("%.4f", v2)
	case bool:
		return fmt.Sprintf("%t", v1) == fmt.Sprintf("%t", v2)
	default:
		return NormalizeString(fmt.Sprintf("%v", v1)) == NormalizeString(fmt.Sprintf("%v", v2))
	}
}
