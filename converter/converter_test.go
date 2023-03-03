package converter_test

import (
	"erlog/converter"
	"erlog/models"
	"testing"
)

func TestConvertBool(t *testing.T) {
	erlog := models.ErLog{BoolKeys: []string{"test", "test2"}, BoolValues: []bool{true, false}}
	converter.ConvertBool(erlog)
}