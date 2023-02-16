package parser

import (
	"errors"
	"fmt"

	"github.com/valyala/fastjson"
)


func ParseJson(value *fastjson.Value) error {
	key := ""
	ParseValue(value, &key)

	return nil
}

func ParseValue(value *fastjson.Value, key *string) error {
	if key == nil {
		return errors.New("key argument is nil")
	}

	switch value.Type() {
	case fastjson.TypeArray:
		val, err := value.Array()

		if err != nil {
			return err
		}

		ParseArray(val, key)

		break
	case fastjson.TypeObject:
		obj, err := value.Object()

		if err != nil {
			return err
		}

		ParseObject(obj, key)
		break
	case fastjson.TypeString:
		// use MarshalTo for speed
		fmt.Printf("Type: str, key: %s, val: %s\n", *key, value.String())
		break
	case fastjson.TypeNumber:
		num, err := value.Float64()

		if err != nil {
			return err
		}

		fmt.Printf("Type: number, key: %s, val: %v\n", *key, num)
		break
	case fastjson.TypeTrue:
		b, err := value.Bool()

		if err != nil {
			return err
		}

		fmt.Printf("Type: bool, key: %s, val: %t\n", *key, b)
		break
	case fastjson.TypeFalse:
		b, err := value.Bool()

		if err != nil {
			return err
		}

		fmt.Printf("Type: bool, key: %s, val: %t\n", *key, b)
		break
	case fastjson.TypeNull:
		// we don't save null types in the database
		break
	}

	return nil
}

func ParseObject(value *fastjson.Object, key *string) error {
	if key == nil {
		return errors.New("key argument is nil")
	}

	value.Visit(func(k []byte, v *fastjson.Value) {
		var value_key string

		if *key == "" {
			value_key = string(k)
		} else {
			value_key = *key + "." + string(k)
		}

		// ok so this is the key and we need to keep track of it
		ParseValue(v, &value_key)
		// basically just check if type is array or nested object
		// for anything else just append that to the specific field array
	})

	return nil
}

func ParseArray(value []*fastjson.Value, key *string) error {
	if key == nil {
		return errors.New("key argument is nil")
	}

	// todo: we need to edit the key or something for the array and figure out how to store it in clickhouse tables
	for _, val := range value {
		ParseValue(val, key)
	}

	// will consume array and return an array of fastjson valuesf
	// if value.Type() != fastjson.TypeArray {
	// 	return fastjson.Value{}, errors.New("Type is not array")
	// }

	return nil
}