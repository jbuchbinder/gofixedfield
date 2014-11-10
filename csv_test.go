package gofixedfield

import (
	"testing"
)

const (
	csvBasicParseTestString   = "1,2,ABC,XYZ"
	csvLayeredParseTestString = "2009-10-10,EX"
)

type csvBasicParseTest struct {
	NumberA int    `csv:"1"`
	NumberB int    `csv:"2"`
	StringC string `csv:"3"`
	StringD string `csv:"4"`
	RawLine string `csv:"raw"`
}

type csvLayeredParseTest struct {
	DateField   *csvDateStruct `csv:"1" csvsplit:"-"`
	StringAfter string         `csv:"2"`
	RawLine     string         `csv:"raw"`
}

type csvDateStruct struct {
	Y int `csv:"1"`
	M int `csv:"2"`
	D int `csv:"3"`
}

func TestCsvBasicParsing(t *testing.T) {
	t.Log("CSV Basic parsing test")
	var out csvBasicParseTest
	UnmarshalCsv(csvBasicParseTestString, ",", &out)
	if out.NumberA != 1 {
		t.Errorf("NumberA parsed as %d", out.NumberA)
	}
	if out.NumberB != 2 {
		t.Errorf("NumberB parsed as %d", out.NumberB)
	}
	if out.StringC != "ABC" {
		t.Errorf("StringC parsed as '%s'", out.StringC)
	}
	if out.StringD != "XYZ" {
		t.Errorf("StringD parsed as '%s'", out.StringD)
	}
	if out.RawLine != csvBasicParseTestString {
		t.Errorf("RawLine parsed as '%s'", out.RawLine)
	}
}

func TestCsvLayeredParsing(t *testing.T) {
	t.Log("CSV Layered parsing test")
	var out csvLayeredParseTest
	UnmarshalCsv(csvLayeredParseTestString, ",", &out)
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
	if out.RawLine != csvLayeredParseTestString {
		t.Errorf("RawLine parsed as '%s'", out.RawLine)
	}
}
