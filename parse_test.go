package gofixedfield

import (
	"testing"
)

const (
	BASIC_PARSE_TEST   = "1234567890ABCDEFGHIJ"
	LAYERED_PARSE_TEST = "20091010EX"
)

type basicParseTest struct {
	NumberA int    `fixed:"1-5"`
	NumberB int    `fixed:"3-5"` // test overlap
	StringC string `fixed:"11-15"`
	StringD string `fixed:"30-35"` // should fail
}

type layeredParseTest struct {
	DateField   *dateStruct `fixed:"1-8"`
	StringAfter string      `fixed:"9-10"`
}

type dateStruct struct {
	Y int `fixed:"1-4"`
	M int `fixed:"5-6"`
	D int `fixed:"7-8"`
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

func TestLayeredParsing(t *testing.T) {
	t.Log("Layered parsing test")
	var out layeredParseTest
	Unmarshal(LAYERED_PARSE_TEST, &out)
	if out.StringAfter != "EX" {
		t.Errorf("Failed to parse after embedded struct/ptr\n")
	}
	if out.DateField.Y != 2009 {
		t.Errorf("Failed to parse embedded Y (Y=%d)\n", out.DateField.Y)
	}
	if out.DateField.M != 10 {
		t.Errorf("Failed to parse embedded M (M=%d)\n", out.DateField.M)
	}
	if out.DateField.D != 10 {
		t.Errorf("Failed to parse embedded D (D=%d)\n", out.DateField.D)
	}
}
