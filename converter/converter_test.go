package converter_test

import (
	"erlog/converter"
	"erlog/models"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestConvertBool(t *testing.T) {
	obj := fastjson.Object{}
	converter := converter.New()

	erlog := models.ErLog{BoolKeys: []string{"test", "test", "test2", "test"}, BoolValues: []bool{true, false, false, true}}
	converter.ConvertBool(erlog, &obj)

	fmt.Printf("%v\n", obj.String())
	assert.Equal(t, obj.String(), "{\"test\":[true,false,true],\"test2\":false}")
}

func TestConvertNumber(t *testing.T) {
	obj := fastjson.Object{}
	converter := converter.New()

	erlog := models.ErLog{NumberKeys: []string{"test", "test", "test2", "test"}, NumberValues: []float64{2, 3.2, 4, 5}}
	converter.ConvertFloat(erlog, &obj)

	fmt.Printf("%v\n", obj.String())
	assert.Equal(t, obj.String(), "{\"test\":[2,3.2,5],\"test2\":4}")
}

func TestConvertString(t *testing.T) {
	obj := fastjson.Object{}
	converter := converter.New()

	erlog := models.ErLog{StringKeys: []string{"test", "test", "test2", "test"}, StringValues: []string{"hi", "hi2", "bye", "hi3"}}
	converter.ConvertString(erlog, &obj)

	fmt.Printf("%v\n", obj.String())
	assert.Equal(t, obj.String(), "{\"test\":[\"hi\",\"hi2\",\"hi3\"],\"test2\":\"bye\"}")
}

func TestConvert(t *testing.T) {
	converter := converter.New()
	uid, err := uuid.Parse("7596cc99-25b3-476e-9b92-4584be4b6478")

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		t.Fail()
	}

	erlog := models.ErLog{
		Id: uid,
		Timestamp: 123123,
		ServiceName: "test_service",
		StringKeys: []string{"testStr", "testStr", "test2Str", "testStr"},
		StringValues: []string{"hi", "hi2", "bye", "hi3"},
		NumberKeys: []string{"testNum", "testNum", "test2Num", "testNum"},
		NumberValues: []float64{2, 3.2, 4, 5},
		BoolKeys: []string{"testBool", "testBool", "test2Bool", "testBool"},
		BoolValues: []bool{true, false, false, true},
	}

	obj, err := converter.Convert(erlog)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		t.Fail()
	}

	fmt.Printf("%v\n", obj.String())
	assert.Equal(t, obj.String(),"{\"id\":\"7596cc99-25b3-476e-9b92-4584be4b6478\",\"serviceName\":\"test_service\",\"timestamp\":123123,\"testStr\":[\"hi\",\"hi2\",\"hi3\"],\"test2Str\":\"bye\",\"testNum\":[2,3.2,5],\"test2Num\":4,\"testBool\":[true,false,true],\"test2Bool\":false}")
}
