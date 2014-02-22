package gofixedfield

import (
	"testing"
)

const (
	BASIC_PARSE_TEST = "1234567890ABCDEFGHIJ"
)

type basicParseTest struct {
	NumberA int    `fixed:"1-5"`
	NumberB int    `fixed:"3-5"` // test overlap
	StringC string `fixed:"11-15"`
	StringD string `fixed:"30-35"` // should fail
}

func TestBasicParsing(t *testing.T) {
	t.Log("Basic parsing test")
	var out basicParseTest
	Unmarshal(BASIC_PARSE_TEST, &out)
	if out.NumberA != 12345 {
		t.Errorf("NumberA parsed as %d", out.NumberA)
	}
	if out.NumberB != 345 {
		t.Errorf("NumberB parsed as %d", out.NumberB)
	}
	if out.StringC != "ABCDE" {
		t.Errorf("StringC parsed as '%s'", out.StringC)
	}
	if out.StringD != "" {
		t.Errorf("StringD should have failed to parse")
	}
}
