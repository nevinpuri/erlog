package converter

import (
	"erlog/models"
	"fmt"

	"github.com/valyala/fastjson"
)

func Convert(erlog models.ErLog) error {
	// todo: use arenapool
	arena := fastjson.Arena{}
	obj := fastjson.Object{}

	// for i, key := range erlog.BoolKeys {
	// 	var val *fastjson.Value

	// 	switch erlog.BoolValues[i] {
	// 	case true:
	// 		val = arena.NewTrue()
	// 		break
	// 	case false:
	// 		val = arena.NewFalse()
	// 		break
	// 	}

	// 	if cval := obj.Get(key); cval != nil {
	// 		switch cval.Type() {
	// 		case fastjson.TypeArray:
	// 			arr, err := cval.Array()

	// 			if err != nil {
	// 				return err
	// 			}

	// 			length := len(arr)

	// 			cval.SetArrayItem(length, val)

	// 			break
	// 		case fastjson.TypeTrue, fastjson.TypeFalse:
	// 			break
	// 		}

	// 		arr := arena.NewArray()
	// 		// arr, err := cval.Array()

	// 		arr = append(arr, val)

	// 		obj.Set(key, arr)

	// 		// this is an array
	// 	} else {
	// 		obj.Set(key, val)
	// 	}

	// }

	for i, key := range erlog.NumberKeys {
		// probably check if the number can be an int, make it an int
		val := arena.NewNumberFloat64(erlog.NumberValues[i])

		obj.Set(key, val)
	}

	for i, key := range erlog.StringKeys {
		// probably check if the number can be an int, make it an int
		val := arena.NewNumberFloat64(erlog.NumberValues[i])

		obj.Set(key, val)
	}


	return nil
}

// todo: add dest
func ConvertBool(erlog models.ErLog) {
	obj := fastjson.Object{}
	arena := fastjson.Arena{}

	for i, key := range erlog.BoolKeys {
		var val *fastjson.Value

		switch erlog.BoolValues[i] {
		case true:
			val = arena.NewTrue()
		case false:
			val = arena.NewFalse()
		}

		obj.Set(key, val)
	}

	fmt.Printf("%v\n", obj.String())
}