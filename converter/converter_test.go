package converter_test

import (
	"erlog/converter"
	"erlog/models"
	"testing"
)

func TestConvertBool(t *testing.T) {
	erlog := models.ErLog{BoolKeys: []string{"test", "test", "test2", "test"}, BoolValues: []bool{true, false, false, true}}
	converter.ConvertBool(erlog)
}