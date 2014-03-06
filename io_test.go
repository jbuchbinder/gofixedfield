package gofixedfield

import (
	"testing"
)

func TestRecordsFromFile(t *testing.T) {
	s, err := RecordsFromFile("./test.txt", EOL_UNIX)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(s) != 3 {
		t.Errorf("Deserialized with %d records\n", len(s))
	}
	if s[0] != "123" || s[1] != "ABC" {
		t.Errorf("Failed to deserialize properly\n")
	}
}
