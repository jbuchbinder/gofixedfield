package gofixedfield

import (
	"fmt"
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
// 		ValA string        `fixed:"1-5"`
//		ValB int           `fixed:"10-15"`
//		ValC *EmbeddedType `fixed:"16-22"`
// 	}
//	type EmbeddedType struct {
//		ValX string `fixed:"1-3"`
//		ValY string `fixed:"4-6"`
//	}
//
//	var out SomeType
//	err := Unmarshal("some string here", &out)
//
// String offsets are one based, not zero based.
func Unmarshal(data string, v interface{}) error {
	//debugStruct(v)
	var val reflect.Value
	if reflect.TypeOf(v).Name() != "" {
		val = reflect.ValueOf(v)
	} else {
		val = reflect.ValueOf(v).Elem()
	}

	//fmt.Printf("Found %d fields\n", val.NumField())
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		cRange := tag.Get("fixed")
		cBookend := strings.Split(cRange, "-")
		if len(cBookend) != 2 {
			// If we don't have two values, skip
			//fmt.Println("Two tag values not found")
			continue
		}

		b, _ := strconv.Atoi(cBookend[0])
		e, _ := strconv.Atoi(cBookend[1])

		b -= 1
		//e -= 1

		// Sanity check range before dying miserably
		if b < 0 || e > len(data) {
			//fmt.Printf("Failed sanity check for b = %d, e = %d, len(data) = %d\n", b, e, len(data))
			continue
		}

		s := data[b:e]

		//fmt.Printf("Field found of type %s\n", typeField.Type.Kind())

		switch typeField.Type.Kind() {
		case reflect.String:
			//fmt.Printf("Found string value '%s'\n", s)
			val.Field(i).SetString(s)
			break
		case reflect.Int:
			//fmt.Printf("Found value '%s'\n", s)
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetInt(v)
			break
		case reflect.Uint:
			//fmt.Printf("Found uint value '%s'\n", s)
			v, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetUint(v)
			break
		case reflect.Ptr, reflect.Struct:
			//fmt.Printf("Found ptr/str value '%s'\n", s)

			// Handle embedded objects by recursively parsing
			// the object with the range we passed.
			if val.Field(i).IsNil() {
				// Initialize pointer to avoid panic
				val.Field(i).Set(reflect.New(val.Field(i).Type().Elem()))
			}
			err := Unmarshal(s, val.Field(i).Interface())
			if err != nil {
				//fmt.Println(err.Error())
			}
			break
		default:
			//fmt.Println("Found unknown value '%s'", s)
			break
		}
	}
	return nil
}

func debugStruct(i interface{}) {
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		fmt.Printf("Field Name: %s, Field Value: %v, Tag Value: %s, Type: %s\n", typeField.Name, valueField.Interface(), tag.Get("fixed"), typeField.Type.Kind())
	}
}
