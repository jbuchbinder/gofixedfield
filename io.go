package gofixedfield

import (
	"io/ioutil"
	"strings"
)

const (
	// EOL_UNIX represents Unix/Linux style end of line.
	EOL_UNIX = "\n"
	// EOL_MAC represents Macintosh style end of line.
	EOL_MAC = "\r"
	// EOL_DOS represents DOS/Windows style end of line.
	EOL_DOS = "\r\n"
)

// RecordsFromFile reads a file and splits into single line records, which
// can be unmarshalled.
func RecordsFromFile(filename string, eolstyle string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), eolstyle), nil
}
