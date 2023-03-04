package converter

import (
	"erlog/models"

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
func ConvertBool(erlog models.ErLog, obj *fastjson.Object) {
	arena := fastjson.Arena{}

	for i, key := range erlog.BoolKeys {
		var val *fastjson.Value

		switch erlog.BoolValues[i] {
		case true:
			val = arena.NewTrue()
		case false:
			val = arena.NewFalse()
		}

		if cval := obj.Get(key); cval != nil {
			// then an object exists
			switch cval.Type() {
			case fastjson.TypeTrue, fastjson.TypeFalse:
				arr := arena.NewArray()

				// first we append the current val
				arr.SetArrayItem(0, cval)
				// then the new value found
				arr.SetArrayItem(1, val)
				obj.Set(key, arr)
			case  fastjson.TypeArray:
				arr, _ := cval.Array()
				l := len(arr)
				cval.SetArrayItem(l, val)
				obj.Set(key, cval)
			}
		} else {
			obj.Set(key, val)
		}
	}

	// fmt.Printf("%v\n", obj.String())
}