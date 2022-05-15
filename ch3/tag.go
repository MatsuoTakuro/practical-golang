package ch3

import (
	"log"
	"reflect"
	"strconv"
	"unsafe"
)

type MapStruct struct {
	Str     string  `map:"str"`  // t.Field(0)
	StrPtr  *string `map:"str"`  // t.Field(1)
	Bool    bool    `map:"bool"` // t.Field(2)
	BoolPtr *bool   `map:"bool"` // t.Field(3)
	Int     int     `map:"int"`  // t.Field(4)
	IntPtr  *int    `map:"int"`  // t.Field(5)
}

func tag() {
	src := map[string]string{
		"str":  "string data",
		"bool": "true",
		"int":  "12345",
	}
	var ms MapStruct
	Decode(&ms, src)
	log.Printf("%+v(%T)\n", ms, ms)
	log.Printf("%+v(%T)\n", ms.Str, ms.Str)
	log.Printf("%+v(%T)\n", ms.StrPtr, ms.StrPtr)
	log.Printf("%+v(%T)\n", ms.Bool, ms.Bool)
	log.Printf("%+v(%T)\n", ms.BoolPtr, ms.BoolPtr)
	log.Printf("%+v(%T)\n", ms.Int, ms.Int)
	log.Printf("%+v(%T)\n", ms.IntPtr, ms.IntPtr)

	src2 := MapStruct{
		Str:     "string-value",
		StrPtr:  &[]string{"string-ptr-value"}[0],
		Bool:    true,
		BoolPtr: &[]bool{true}[0],
		Int:     12345,
		IntPtr:  &[]int{67890}[0],
	}
	dest := map[string]string{}
	Encode(dest, &src2)
	for k, v := range dest {
		log.Printf("%s(%T) : %s(%T)\n", k, k, v, v)
	}
}

func Decode(target interface{}, src map[string]string) error {
	v := reflect.ValueOf(target)
	e := v.Elem()
	return decode(e, src)
}

func decode(e reflect.Value, src map[string]string) error {
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Anonymous {
			if err := decode(e.Field(i), src); err != nil {
				return err
			}
			continue
		}

		if f.Type.Kind() == reflect.Struct {
			if err := decode(e.Field(i), src); err != nil {
				return err
			}
			continue
		}

		key := f.Tag.Get("map")
		if key == "" {
			key = f.Name
		}

		v, ok := src[key]
		if !ok {
			continue
		}

		var k reflect.Kind
		var isP bool
		if f.Type.Kind() != reflect.Ptr {
			k = f.Type.Kind()
		} else {
			k = f.Type.Elem().Kind()
			if k == reflect.Ptr {
				continue
			}
			isP = true
		}

		switch k {
		case reflect.String:
			if isP {
				e.Field(i).Set(reflect.ValueOf(&v))
			} else {
				e.Field(i).SetString(v)
			}
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err == nil {
				if isP {
					e.Field(i).Set(reflect.ValueOf(&b))
				} else {
					e.Field(i).SetBool(b)
				}
			}
		case reflect.Int:
			n64, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				if isP {
					n := int(n64)
					e.Field(i).Set(reflect.ValueOf(&n))
				} else {
					e.Field(i).SetInt(n64)
				}
			}
		}
	}
	return nil
}

func Encode(target map[string]string, src interface{}) error {
	v := reflect.ValueOf(src)
	e := v.Elem()
	t := e.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Anonymous {
			if err := Encode(target, e.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if f.Type.Kind() == reflect.Struct {
			if err := Encode(target, e.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		key := f.Tag.Get("map")
		if key == "" {
			key = f.Name
		}

		var k reflect.Kind
		var isP bool
		if f.Type.Kind() != reflect.Ptr {
			k = f.Type.Kind()
		} else {
			k = f.Type.Elem().Kind()
			if k == reflect.Ptr {
				continue
			}
			isP = true
		}

		switch k {
		case reflect.String:
			if isP {
				if e.Field(i).Pointer() != 0 {
					target[key] = *(*string)(unsafe.Pointer(e.Field(i).Pointer()))
				}
			} else {
				target[key] = e.Field(i).String()
			}
		case reflect.Bool:
			var b bool
			if isP {
				if e.Field(i).Pointer() != 0 {
					b = *(*bool)(unsafe.Pointer(e.Field(i).Pointer()))
				}
			} else {
				b = e.Field(i).Bool()
			}
			target[key] = strconv.FormatBool(b)
		case reflect.Int:
			var n int64
			if isP {
				if e.Field(i).Pointer() != 0 {
					n = int64(*(*int)(unsafe.Pointer(e.Field(i).Pointer())))
				}
			} else {
				n = e.Field(i).Int()
			}
			target[key] = strconv.FormatInt(n, 10)
		}
	}
	return nil
}
