package gofixedfield

import (
	"os"
	"strings"
)

const (
	// EOLUnix represents Unix/Linux style end of line.
	EOLUnix = "\n"
	// EOLMac represents Macintosh style end of line.
	EOLMac = "\r"
	// EOLDOS represents DOS/Windows style end of line.
	EOLDOS = "\r\n"
)

// DecimalComma enables the parsing of numeric values having a comma
// instead of a point as decimal separator.
var DecimalComma bool

// RecordsFromFile reads a file and splits into single line records, which
// can be unmarshalled.
func RecordsFromFile(filename string, eolstyle string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), eolstyle), nil
}
