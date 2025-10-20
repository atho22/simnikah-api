package config

import "gorm.io/gorm"

// InDB represents a database handler that can be embedded in other structs
type InDB struct {
	DB *gorm.DB
}
