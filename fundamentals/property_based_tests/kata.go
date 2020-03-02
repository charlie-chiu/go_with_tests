package property_based_tests

import (
	"strings"
)

type RomanNumeral struct {
	Value  int
	Symbol string
}

var allRomanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func ConvertToRoman(arabic int) string {
	var result strings.Builder

	for _, romanNumeral := range allRomanNumerals {
		for arabic >= romanNumeral.Value {
			result.WriteString(romanNumeral.Symbol)
			arabic -= romanNumeral.Value
		}
	}

	return result.String()
}

func convertToArabic(roman string) (arabic int) {
	//todo: complete this
	return arabic
}
