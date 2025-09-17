package data

import (
	"strconv"
	"testing"
)

func TestIsFloat_CompareWithParseFloat(t *testing.T) {
	testStrings := []string{
		"123.456", "123", ".456", "123.", "+123.456", "-123.456",
		"1.23e10", "1.23E10", "1.23e-10", "1.23E+10", "0", "0.0",
		"-0", "+0", "Inf", "-Inf", "+Inf", "inf", "-inf", "+inf",
		"NaN", "nan", "", "abc", "12.34.56", "12e", "12e+",
		"12ee10", ".e10", "e10", "12..34", "++12", "--12",
		"12-34", "12+34", " 123 ", "∞",
	}

	for _, s := range testStrings {
		_, parseErr := strconv.ParseFloat(s, 64)
		canParse := parseErr == nil
		regexResult := IsFloat(s)

		if canParse != regexResult {
			t.Errorf("Mismatch for %q: ParseFloat=%v, IsFloat=%v", s, canParse, regexResult)
		}
	}
}
