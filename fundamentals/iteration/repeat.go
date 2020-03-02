package iteration

import "strings"

// Repeat : repeat a character 5 times
func Repeat(character string, times int) string {
	var b strings.Builder

	for i := 0; i < times; i++ {
		b.WriteString(character)
	}

	return b.String()
}
