package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Title    string
	AudioKey string
}
