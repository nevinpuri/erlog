package parser

import (
	"erlog/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)


func ParseJson(value *fastjson.Value) (models.ErLog, error) {
	key := ""
	erlog := models.ErLog{}
	erlog.Id = uuid.New()

	timestamp := value.GetInt("timestamp")

	if timestamp != 0 {
		erlog.Timestamp = int64(timestamp)
	}

	err := ParseValue(value, &key, &erlog)

	if err != nil {
		return models.ErLog{}, err
	}

	return erlog, nil
}

func ParseValue(value *fastjson.Value, key *string, erlog *models.ErLog) error {
	if key == nil {
		return errors.New("key argument is nil")
	}
	
	if erlog == nil {
		return errors.New("Erlog argument is nil")
	}

	switch value.Type() {
	case fastjson.TypeArray:
		val, err := value.Array()

		if err != nil {
			return err
		}

		ParseArray(val, key, erlog)

		break
	case fastjson.TypeObject:
		obj, err := value.Object()

		if err != nil {
			return err
		}

		ParseObject(obj, key, erlog)
		break
	case fastjson.TypeString:
		// use MarshalTo for speed

		// no idea what marshalto does except for its faster than string
		dst, err := value.StringBytes()

		if err != nil {
			return err
		}

		erlog.StringKeys = append(erlog.StringKeys, *key)
		erlog.StringValues = append(erlog.StringValues, string(dst))

		fmt.Printf("Type: str, key: %s, val: %s\n", *key, string(dst))
		break
	case fastjson.TypeNumber:
		if *key == "timestamp" {
			break
		}

		num, err := value.Float64()

		if err != nil {
			return err
		}

		erlog.NumberKeys = append(erlog.NumberKeys, *key)
		erlog.NumberValues = append(erlog.NumberValues, num)

		fmt.Printf("Type: number, key: %s, val: %v\n", *key, num)
		break
	case fastjson.TypeTrue:
		b, err := value.Bool()

		if err != nil {
			return err
		}

		erlog.BoolKeys = append(erlog.BoolKeys, *key)
		erlog.BoolValues = append(erlog.BoolValues, b)

		fmt.Printf("Type: bool, key: %s, val: %t\n", *key, b)
		break
	case fastjson.TypeFalse:
		b, err := value.Bool()

		if err != nil {
			return err
		}
		
		erlog.BoolKeys = append(erlog.BoolKeys, *key)
		erlog.BoolValues = append(erlog.BoolValues, b)

		fmt.Printf("Type: bool, key: %s, val: %t\n", *key, b)
		break
	case fastjson.TypeNull:
		// we don't save null types in the database
		break
	}

	return nil
}

func ParseObject(value *fastjson.Object, key *string, erlog *models.ErLog) error {
	if key == nil {
		return errors.New("key argument is nil")
	}
	
	if erlog == nil {
		return errors.New("Erlog argument is nil")
	}

	value.Visit(func(k []byte, v *fastjson.Value) {
		var value_key string

		if *key == "" {
			value_key = string(k)
		} else {
			value_key = *key + "." + string(k)
		}

		// ok so this is the key and we need to keep track of it
		ParseValue(v, &value_key, erlog)
		// basically just check if type is array or nested object
		// for anything else just append that to the specific field array
	})

	return nil
}

func ParseArray(value []*fastjson.Value, key *string, erlog *models.ErLog) error {
	if key == nil {
		return errors.New("key argument is nil")
	}
	
	if erlog == nil {
		return errors.New("Erlog argument is nil")
	}

	// todo: we need to edit the key or something for the array and figure out how to store it in clickhouse tables
	for _, val := range value {
		ParseValue(val, key, erlog)
	}

	// will consume array and return an array of fastjson valuesf
	// if value.Type() != fastjson.TypeArray {
	// 	return fastjson.Value{}, errors.New("Type is not array")
	// }

	return nil
}