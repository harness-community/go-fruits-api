package utils

import (
	"os"
	"strings"
)

//Reverse reverses the String
func Reverse(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i := len(runes) - 1; 0 <= i; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}

//LookupEnvOrString looks up an environment variable if not found
//returns defaultVal
func LookupEnvOrString(envName, defaultVal string) string {
	if val, ok := os.LookupEnv(envName); ok {
		return val
	}

	return defaultVal
}
