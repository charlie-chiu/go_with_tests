package property_based_tests

import (
	"fmt"
	"testing"
)

var cases = []struct {
	Arabic int
	Roman  string
}{
	{1, "I"},
	{2, "II"},
	{3, "III"},
	{4, "IV"},
	{5, "V"},
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
}

func TestRomanNumerals(t *testing.T) {
	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%d gets converted to %q", testCase.Arabic, testCase.Roman), func(t *testing.T) {
			got := ConvertToRoman(testCase.Arabic)
			want := testCase.Roman

			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

//amazing, re-use first test cases!
func TestConvertingToArabic(t *testing.T) {
	for _, test := range cases[:4] {
		t.Run(fmt.Sprintf("%q gots converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
			got := convertToArabic(test.Roman)
			if got != test.Arabic {
				t.Errorf("got %v, want %d", got, test.Arabic)
			}
		})
	}
}
