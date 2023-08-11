package editor

import (
	"strings"
)

func createAddressFromString(address string) []string {
	var separator byte = '.'
	var escapeString byte = '\\'

	var result []string
	var token []byte
	for i := 0; i < len(address); i++ {
		if address[i] == separator {
			result = append(result, string(token))
			token = token[:0]
		} else if address[i] == escapeString && i+1 < len(address) {
			i++
			token = append(token, address[i])
		} else {
			token = append(token, address[i])
		}
	}
	result = append(result, string(token))
	return result
}

func createStringFromAddress(address []string) string {
	separator := "."
	escapeString := "\\"

	result := ""

	for i, s := range address {
		if i > 0 {
			result = result + separator
		}
		result = result + strings.ReplaceAll(s, separator, escapeString+separator)
	}

	return result
}
