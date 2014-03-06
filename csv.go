package gofixedfield

import (
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalCsv unmarshals string data into an annotated interface. This
// should resemble:
//
// 	type SomeType struct {
// 		ValA string        `csv:"1"`
//		ValB int           `csv:"2"`
//		ValC *EmbeddedType `csv:"3" csvsplit:"~"`
// 	}
//	type EmbeddedType struct {
//		ValX string `csv:"1"`
//		ValY string `csv:"2"`
//	}
//
//	var out SomeType
//	err := Unmarshal("A,2,X~Y", "," &out)
//
// String offsets are one based, not zero based.
func UnmarshalCsv(data string, sep string, v interface{}) error {
	//debugStruct(v)
	var val reflect.Value
	if reflect.TypeOf(v).Name() != "" {
		val = reflect.ValueOf(v)
	} else {
		val = reflect.ValueOf(v).Elem()
	}

	//fmt.Println("UnmarshalCsv called with separator " + sep)
	parts := strings.Split(data, sep)

	//fmt.Printf("Found %d fields\n", val.NumField())
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		cField := tag.Get("csv")
		cSep := tag.Get("csvsplit")
		if len(cField) < 1 || len(cField) > 4 {
			//fmt.Println("Bailing out, invalid csv tag ", cField)
			continue
		}

		f, _ := strconv.Atoi(cField)
		f -= 1

		// Sanity check range before dying miserably
		if f < 0 || f > len(parts) {
			//fmt.Printf("Failed sanity check for f = %d, len(parts) = %d\n", f, len(parts))
			continue
		}

		s := parts[f]
		//fmt.Printf("s == %s\n", s)

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
			if cSep == "" {
				//fmt.Println("No csvsplit defined")
				continue
			}
			//fmt.Printf("Found ptr/str value '%s'\n", s)

			// Handle embedded objects by recursively parsing
			// the object with the range we passed.
			if val.Field(i).IsNil() {
				// Initialize pointer to avoid panic
				val.Field(i).Set(reflect.New(val.Field(i).Type().Elem()))
			}
			err := UnmarshalCsv(s, cSep, val.Field(i).Interface())
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
