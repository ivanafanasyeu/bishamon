package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input    float64
	expected string
}

func TestFormatNumberToCurrency(t *testing.T) {
	testCases := []testCase{
		{input: 100, expected: "100.00"},
		{input: 100.00, expected: "100.00"},
		{input: 100.00000, expected: "100.00"},
		{input: 1230, expected: "1,230.00"},
		{input: 888123, expected: "888,123.00"},
		{input: 0, expected: "0.00"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("FormatNumberToCurrency(%v)", tc.input), func(t *testing.T) {
			output := FormatNumberToCurrency(tc.input)
			assert.Equal(t, tc.expected, output)
		})
	}
}
