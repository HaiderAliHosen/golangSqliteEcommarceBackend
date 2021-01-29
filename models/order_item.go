package models

import "github.com/jinzhu/gorm"

//OrderItem __
type OrderItem struct {
	gorm.Model
	Order   Order
	OrderID uint `gorm:"not null"`

	Product   Product
	ProductID uint `gorm:"not null"`

	Slug        string `gorm:"not null"`
	ProductName string `gorm:"not null"`
	Price       int    `gorm:"not null"`
	Quantity    int    `gorm:"not null"`

	User   User `gorm:"association_foreignkey:UserId:"`
	UserID uint `gorm:"default:null"`
}
