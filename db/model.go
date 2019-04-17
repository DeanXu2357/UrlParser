package db

import (
	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
)

type AccessLog struct {
	gorm.Model
	// ID          uint `gorm:"primary_key;AUTO_INCREMENT"`
	IP          string
	Method      string
	Route       string
	Domain      string
	TriggeredAt string
}

type AccessLogParam struct {
	gorm.Model
	// ID          uint `gorm:"primary_key;AUTO_INCREMENT"`
	AccessLogID uint
	Key         string
	Value       string
}
