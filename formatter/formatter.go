package formatter

import (
	"strings"

	"github.com/valyala/fastjson"
)

func FormatObj(obj fastjson.Object) (*fastjson.Object, error) {
	arena := fastjson.Arena{}
	out, err := arena.NewObject().Object()

	if err != nil {
		return out, err
	}

	obj.Visit(func(key []byte, v *fastjson.Value) {
		children := strings.Split(string(key), ".")
		keyVal := children[len(children)-1]

		// remove last element since we don't want to create it as an object
		children = children[:len(children)-1]
		subObj, cerr := CreateSubObjects(out, children)

		if cerr != nil {
			err = cerr
			return
		}

		subObj.Set(keyVal, v)
	})

	return out, err
}

// should both create all the objects, and also merge any which already have values
// we assume the json is validated already
func CreateSubObjects(obj *fastjson.Object, children []string) (*fastjson.Object, error) {
	// just make this an arena pool or make it passed by reference

	arena := fastjson.Arena{}
	childObj := obj

	var err error

	for _, key := range children {
		current := childObj.Get(key)

		if current == nil {
			newObj := arena.NewObject()
			childObj.Set(key, newObj)
		}

		childObj, err = childObj.Get(key).Object()

		if err != nil {
			return childObj, err
		}
	}

	return childObj, err
}