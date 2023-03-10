package models

import (
	"github.com/google/uuid"
)

type ErLog struct {
	Id				uuid.UUID		`json:"id" ch:"Id"`
	Timestamp		int64		`json:"timestamp" ch:"Timestamp"`
	ServiceName		string			`json:"serviceName" ch:"ServiceName"`
	StringKeys		[]string		`json:"stringKeys" ch:"StringKeys"`
	StringValues	[]string		`json:"stringValues" ch:"StringValues"`
	NumberKeys		[]string		`json:"numberKeys" ch:"NumberKeys"`
	NumberValues	[]float64		`json:"numberValues" ch:"NumberValues"`
	BoolKeys		[]string		`json:"boolKeys" ch:"BoolKeys"`
	BoolValues		[]bool			`json:"boolValues" ch:"BoolValues"`
	Raw				string			`json:"raw"`
}

// why don't we get the clickhouse tables going first and then just use go to infer the types in code