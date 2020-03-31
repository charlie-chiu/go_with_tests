package kata_test

import (
	"fmt"
	kata "github.com/charlie-chiu/go_with_test/fundamentals/property_based_tests"
	"testing"
	"testing/quick"
)

var cases = []struct {
	Arabic uint16
	Roman  string
}{
	{1, "I"},
	{2, "II"},
	{3, "III"},
	{4, "IV"},
	{5, "V"},
	{6, "VI"},
	{8, "VIII"},
	{9, "IX"},
	{10, "X"},
	{14, "XIV"},
	{18, "XVIII"},
	{20, "XX"},
	{29, "XXIX"},
	{39, "XXXIX"},
	{40, "XL"},
	{50, "L"},
	{60, "LX"},
	{68, "LXVIII"},
	{100, "C"},
	{90, "XC"},
	{125, "CXXV"},
	{500, "D"},
	{400, "CD"},
	{888, "DCCCLXXXVIII"},
	{900, "CM"},
	{1000, "M"},
	{1984, "MCMLXXXIV"},
	{1006, "MVI"},
	{3999, "MMMCMXCIX"},
}

func TestRomanNumerals(t *testing.T) {
	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%d gets converted to %q", testCase.Arabic, testCase.Roman), func(t *testing.T) {
			got := kata.PrintInRoman(testCase.Arabic)
			want := testCase.Roman

			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestConvertToArabic(t *testing.T) {
	for _, test := range cases {
		t.Run(fmt.Sprintf("%q converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
			got := kata.ConvertToArabic(test.Roman)
			if got != test.Arabic {
				t.Errorf("got %d want %d", got, test.Arabic)
			}
		})
	}
}

func TestPropertiesOfConversion(t *testing.T) {
	assertion := func(arabic uint16) bool {
		if arabic > 3999 {
			// how to do this?
			// generate arabic within 1 ~ 3999
			// or return an error when out of range ?
			return true
		}
		t.Log("testing", arabic)
		roman := kata.PrintInRoman(arabic)
		fromRoman := kata.ConvertToArabic(roman)
		return fromRoman == arabic
	}

	config := &quick.Config{
		MaxCount: 100,
		Values:   nil,
	}

	if err := quick.Check(assertion, config); err != nil {
		t.Error("failed checks", err)
	}
}
