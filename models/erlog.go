package models

import (
	"time"

	"github.com/google/uuid"
)

type ErLog struct {
	Id				uuid.UUID		`json:"id" ch:"Id"`
	Timestamp		time.Time		`json:"timestamp" ch:"Timestamp"`
	ServiceName		string			`json:"serviceName" ch:"ServiceName"`
	StringKeys		[]string		`json:"stringKeys" ch:"String.Keys"`
	StringValues	[]string		`json:"stringValues" ch:"String.Values"`
	NumberKeys		[]string		`json:"numberKeys" ch:"Number.Keys"`
	NumberValues	[]float64		`json:"numberValues" ch:"Number.Values"`
	BoolKeys		[]string		`json:"boolKeys" ch:"Bool.Keys"`
	BoolValues		[]uint8			`json:"boolValues" ch:"Bool.Values"`
}

// why don't we get the clickhouse tables going first and then just use go to infer the types in code