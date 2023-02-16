package queue_test

import (
	"erlog/queue"
	"testing"

	"github.com/valyala/fastjson"
)

func TestParseVal(t *testing.T) {
	val, err := fastjson.Parse("{\"key\": \"val\", \"item\": [\"first\", \"second\"], \"obj\": {\"val\": \"hi\"}}")

	if err != nil {
		t.Logf("%s", err.Error())
	}

	keys := ""

	queue.ParseValue(val, &keys)
}