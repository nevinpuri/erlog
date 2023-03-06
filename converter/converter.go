package converter

import (
	"erlog/models"

	"github.com/valyala/fastjson"
)

type Converter struct {
	pool fastjson.ArenaPool
}

func New() Converter {
	return Converter {
		pool: fastjson.ArenaPool{},
	}
}

func (c Converter) Convert(erlog models.ErLog) (fastjson.Object, error) {
	// get all the other metadata from this, like timestamp and id
	arena := c.pool.Get()
	obj := fastjson.Object{}

	var err error

	id := arena.NewString(erlog.Id.String())
	obj.Set("id", id)

	serviceName := arena.NewString(erlog.ServiceName)
	obj.Set("serviceName", serviceName)

	timestamp := arena.NewNumberFloat64(float64(erlog.Timestamp))
	obj.Set("timestamp", timestamp)

	err = c.ConvertString(erlog, &obj)

	if err != nil {
		return fastjson.Object{}, err
	}

	err = c.ConvertFloat(erlog, &obj)

	if err != nil {
		return fastjson.Object{}, err
	}

	err = c.ConvertBool(erlog, &obj)

	if err != nil {
		return fastjson.Object{}, err
	}

	return obj, nil
}

// todo: add dest
func (c Converter) ConvertBool(erlog models.ErLog, obj *fastjson.Object) error {
	arena := c.pool.Get()

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
			case fastjson.TypeArray:
				arr, err := cval.Array()

				if err != nil {
					return err
				}

				l := len(arr)
				cval.SetArrayItem(l, val)
				obj.Set(key, cval)
			}
		} else {
			obj.Set(key, val)
		}
	}

	return nil
}

func (c Converter) ConvertFloat(erlog models.ErLog, obj *fastjson.Object) error {
	arena := c.pool.Get()

	for i, key := range erlog.NumberKeys {
		floatVal := erlog.NumberValues[i]
		val := arena.NewNumberFloat64(floatVal)

		if cval := obj.Get(key); cval != nil {
			switch cval.Type() {
			case fastjson.TypeNumber:
				arr := arena.NewArray()

				arr.SetArrayItem(0, cval)
				arr.SetArrayItem(1, val)
				obj.Set(key, arr)
			case fastjson.TypeArray:
				arr, err := cval.Array()

				if err != nil {
					return err
				}

				l := len(arr)
				cval.SetArrayItem(l, val)
				obj.Set(key, cval)
			}
		} else {
			obj.Set(key, val)
		}
	}

	return nil
}

func (c Converter) ConvertString(erlog models.ErLog, obj *fastjson.Object) error {
	arena := c.pool.Get()

	for i, key := range erlog.StringKeys {
		strVal := erlog.StringValues[i]
		val := arena.NewString(strVal)

		if cval := obj.Get(key); cval != nil {
			switch cval.Type() {
			case fastjson.TypeString:
				arr := arena.NewArray()

				arr.SetArrayItem(0, cval)
				arr.SetArrayItem(1, val)
				obj.Set(key, arr)
			case fastjson.TypeArray:
				arr, err := cval.Array()

				if err != nil {
					return err
				}

				l := len(arr)
				cval.SetArrayItem(l, val)
				obj.Set(key, cval)
			}
		} else {
			obj.Set(key, val)
		}
	}

	return nil
}


// returns prefix, rest_of_string (without leading .)
func ConsumeParent(str string) (string, string) {
	var idx int
	var out string
	var prefix string

	for i, char := range str {
		if char == '.' {
			idx = i + 1
			break
		}
	}

	if idx == 0 {
		return "", ""
	}

	out = str[idx:]
	prefix = str[0:idx-1]
	return prefix, out
}