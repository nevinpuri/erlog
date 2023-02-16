package parser_test

import (
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

	parser.ParseValue(val, &keys)
}