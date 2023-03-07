package parser_test

import (
	"erlog/models"
	"erlog/parser"
	"fmt"
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

func TestArrayObj(t *testing.T) {
	val, err := fastjson.Parse("{\r\n  \"Msg\": \"Job Finished\",\r\n  \"User\": {\r\n    \"ID\": \"cx50yz\",\r\n    \"Name\": [\r\n      \"first\",\r\n      \"last\"\r\n    ]\r\n  }\r\n}")

	if err != nil {
		t.Logf("%s", err.Error())
		t.Fail()
	}

	keys := ""
	erlog := models.ErLog{}

	parser.ParseValue(val, &keys, &erlog)

	fmt.Printf("%+v\n", erlog)
}