package models

import (
	"gorm.io/datatypes"
)

type ErLog struct {
	Model
	Data	datatypes.JSON	`json:"data"`
}