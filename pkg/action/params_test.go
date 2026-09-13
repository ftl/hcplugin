package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseParams(t *testing.T) {
	tt := []struct {
		desc     string
		value    string
		expected map[string]string
	}{
		{desc: "empty", value: "", expected: nil},
		{desc: "blank", value: "   ", expected: nil},
		{desc: "one parameter", value: "number=5", expected: map[string]string{"number": "5"}},
		{desc: "two parameters", value: "amount=100,unit=hz", expected: map[string]string{"amount": "100", "unit": "hz"}},
		{desc: "spaces around the pairs", value: " number = 5 , text = cq dx ", expected: map[string]string{"number": "5", "text": "cq dx"}},
		{desc: "empty value", value: "number=", expected: map[string]string{"number": ""}},
		{desc: "value with an equals sign", value: "expression=a=b", expected: map[string]string{"expression": "a=b"}},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err := parseParams(tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseParams_Invalid(t *testing.T) {
	tt := []string{"number", "number=5,amount", "=5", "number=5,"}
	for _, value := range tt {
		t.Run(value, func(t *testing.T) {
			_, err := parseParams(value)
			assert.Error(t, err)
		})
	}
}
