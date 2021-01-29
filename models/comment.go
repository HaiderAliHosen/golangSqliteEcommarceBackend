package models

import (
	"github.com/jinzhu/gorm"
)

//Comment ---
type Comment struct {
	gorm.Model
	Content   string  `gorm:"size:2048"`
	Rating    int     `gorm:"default:null"`
	Product   Product `gorm:"foreignkey:ProductId"`
	ProductID uint    `gorm:"not null"`
	User      User    `gorm:"foreignkey:UserId"`
	UserID    uint    `gorm:"not null"`
}
