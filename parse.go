package gofixedfield

import (
	"reflect"
	"strconv"
	"strings"
)

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

		b--
		//e--

		// Sanity check range before dying miserably
		if b < 0 || e > len(data) {
			//fmt.Printf("Failed sanity check for b = %d, e = %d, len(data) = %d\n", b, e, len(data))
			continue
		}

		s := data[b:e]

		//fmt.Printf("Field found of type %s\n", typeField.Type.Kind())

		switch typeField.Type.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(s)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetBool(v)
			break
		case reflect.Float32:
			if DECIMAL_COMMA {
				s = strings.Replace(s, ",", ".", 1)
			}
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetFloat(v)
			break
		case reflect.Float64:
			if DECIMAL_COMMA {
				s = strings.Replace(s, ",", ".", 1)
			}
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetFloat(v)
			break
		case reflect.String:
			//fmt.Printf("Found string value '%s'\n", s)
			val.Field(i).SetString(s)
			break
		case reflect.Int8:
			//fmt.Printf("Found value '%s'\n", s)
			v, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetInt(v)
			break
		case reflect.Int32:
			//fmt.Printf("Found value '%s'\n", s)
			v, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}
			val.Field(i).SetInt(v)
			break
		case reflect.Int, reflect.Int64:
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
