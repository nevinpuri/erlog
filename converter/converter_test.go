package converter_test

import (
	"erlog/converter"
	"erlog/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestConvertBool(t *testing.T) {
	obj := fastjson.Object{}
	erlog := models.ErLog{BoolKeys: []string{"test", "test", "test2", "test"}, BoolValues: []bool{true, false, false, true}}
	converter.ConvertBool(erlog, &obj)

	assert.Equal(t, obj.String(), "{\"test\":[true,false,true],\"test2\":false}")
	fmt.Printf("%v\n", obj.String())
}