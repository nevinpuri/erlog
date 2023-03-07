package formatter_test

import (
	"erlog/converter"
	"erlog/formatter"
	"erlog/models"
	"fmt"
	"testing"
)

func TestFormatter(t *testing.T) {
	erlog := models.ErLog{
		StringKeys: []string{"test", "user.id", "user.field", "user.field"},
		StringValues: []string{"hi", "abc", "hi", "bye"},
	}

	converter := converter.New()

	obj, err := converter.Convert(erlog)

	if err != nil {
		fmt.Printf("%v\n", err.Error())
		t.Fail()
	}

	outObj, err := formatter.FormatObj(obj)

	fmt.Printf("%v\n", outObj.String())
}