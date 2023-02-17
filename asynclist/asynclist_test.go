package asynclist_test

import (
	"erlog/asynclist"
	"erlog/models"
	"testing"

	"github.com/google/uuid"
)

func TestClearLen(t *testing.T) {
	l := asynclist.New(10)

	l.Append(models.ErLog{Id: uuid.New()})
	l.Append(models.ErLog{Id: uuid.New()})

	l.Clear()

	l.Append(models.ErLog{Id: uuid.New()})
	// all := l.All()

	// assert.Equal(t, len(all), 1)
}