package main

/*import (
	"reflect"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Uniform struct {
	id int32
}

type (u Uniform) Set(typeStr string, values ...interface{}) {
	if len(typeStr) != 2 && len(typeStr) != 3 {
		log.Panic("Invalid type string")
	}

	var expType reflect.Type
	switch typeStr[1] {
	case 'i':
		expType = reflect.Int32
	case 'u':
		expType = reflect.Uint32
	case 'f':
		expType = reflect.Float32
	case 'd':
		expType = reflect.Float64
	}

	var expNb int
	switch typeStr[0] {
	case '1':
		expNb = 1
	case '2':
		expNb = 2
	case '3':
		expNb = 3
	case '4':
		expNb = 4
	default:
		log.Panic("Invalid type string")
	}

	if len(values) != expNb {

	}

	for _, v := range values {
		if reflect.TypeOf(v).ConvertibleTo(expType) {
			log.Panic("Passed values with incompatible types")
		}
	}

	switch len(values) {
	case 1:
		val := values[0]
		switch reflect.TypeOf(val).Kind() {
			case reflect.I
		}
	}
}*/
