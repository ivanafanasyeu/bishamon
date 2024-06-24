package utils

import (
	"fmt"
	"strings"
)

func FormatNumberToCurrency(n float64) string {
	formatted := fmt.Sprintf("%.2f", n)

	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	fractionalPart := parts[1]

	numLength := len(integerPart)
	if n <= 3 {
		return formatted
	}

	var result strings.Builder
	numCommas := (numLength - 1) / 3
	result.Grow(numLength + numCommas + len(fractionalPart) + 1)

	for i, j := numLength%3, 0; j < numLength; j++ {
		if j != 0 && j%3 == i {
			result.WriteByte(',')
		}
		result.WriteByte(integerPart[j])
	}

	result.WriteByte('.')
	result.WriteString(fractionalPart)

	return result.String()
}
