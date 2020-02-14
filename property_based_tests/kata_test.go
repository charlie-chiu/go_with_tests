package property_based_tests

import (
	"fmt"
	"testing"
)

func TestRomanNumerals(t *testing.T) {
	cases := []struct {
		Arabic int
		Want   string
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

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%d gets converted to %q", testCase.Arabic, testCase.Want), func(t *testing.T) {
			got := ConvertToRoman(testCase.Arabic)
			want := testCase.Want

			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
