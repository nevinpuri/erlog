package parser_test

import (
	"erlog/models"
	"erlog/parser"
	"testing"

	"github.com/valyala/fastjson"
)

func TestParseVal(t *testing.T) {
	val, err := fastjson.Parse("{\"key\": \"val\", \"item\": [3, \"second\"], \"obj\": {\"val\": false}}")

	if err != nil {
		t.Logf("%s", err.Error())
	}

	keys := ""
	erlog := models.ErLog{}

	parser.ParseValue(val, &keys, &erlog)

	t.Logf("String keys: %v, string values: %v, number keys: %v, number values: %v, bool keys: %v, bool_values: %v\n", erlog.StringKeys, erlog.StringValues, erlog.NumberKeys, erlog.NumberValues, erlog.BoolKeys, erlog.BoolValues)
}