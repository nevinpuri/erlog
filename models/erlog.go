package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ErLog struct {
	gorm.Model
	O	datatypes.JSON	`json:"o"`
}