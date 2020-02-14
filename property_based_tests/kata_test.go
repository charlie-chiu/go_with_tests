package property_based_tests

import "testing"

func TestRomanNumerals(t *testing.T) {
	cases := []struct {
		Name   string
		Arabic int
		Want   string
	}{
		{"1 gets converted to I", 1, "I"},
		{"2 gets converted to II", 2, "II"},
		{"3 gets converted to III", 3, "III"},
		{"4 gets converted to IV", 4, "IV"},
		{"5 gets converted to V", 5, "V"},
		{"8 gets converted to VIII", 8, "VIII"},
		{"9 gets converted to IX", 9, "IX"},
		{"10 gets converted to X", 10, "X"},
		{"14 gets converted to XIV", 14, "XIV"},
		{"18 gets converted to XVIII", 18, "XVIII"},
		{"20 gets converted to XX", 20, "XX"},
		{"29 gets converted to XXIX", 29, "XXIX"},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := ConvertToRoman(testCase.Arabic)
			want := testCase.Want

			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
