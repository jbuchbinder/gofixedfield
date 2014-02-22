package gofixedfield

import (
	"io/ioutil"
	"reflect"
	"strconv"
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

// Unmarshal unmarshals string data into an annotated interface. This should
// resemble:
//
// 	type SomeType struct {
// 		ValA string `fixed:"1-5"`
//		ValB int    `fixed:"10-15"`
// 	}
//
//	var out SomeType
//	err := Unmarshal("some string here", &out)
//
// String offsets are one based, not zero based.
func Unmarshal(data string, v interface{}) error {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		cRange := tag.Get("fixed")
		cBookend := strings.Split(cRange, "-")
		if len(cBookend) != 2 {
			// If we don't have two values, skip
			continue
		}

		b, _ := strconv.Atoi(cBookend[0])
		e, _ := strconv.Atoi(cBookend[1])

		b -= 1
		//e -= 1

		// Sanity check range before dying miserably
		if b < 0 || e >= len(data) {
			continue
		}

		s := data[b:e]

		switch typeField.Type.Kind() {
		case reflect.String:
			val.Field(i).SetString(s)
			break
		case reflect.Int:
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				continue
			}
			val.Field(i).SetInt(v)
			break
		case reflect.Uint:
			v, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				continue
			}
			val.Field(i).SetUint(v)
			break
		case reflect.Struct:
			// Handle embedded objects by recursively parsing
			// the object with the range we passed.
			Unmarshal(s, val.Field(i).Interface())
			break
		default:
			break
		}
	}
	return nil
}
