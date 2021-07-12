package utils

import "strings"

//Reverse reverses the String
func Reverse(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i := len(runes) - 1; 0 <= i; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}
