package models

import "github.com/jinzhu/gorm"

//FileUpload ---
type FileUpload struct {
	gorm.Model
	Filename     string
	FilePath     string
	OriginalName string
	FileSize     uint

	Tag   Tag  `gorm:"association_foreignkey:TagId"`
	TagID uint `gorm:"default:null"`

	Category   Category `gorm:"association_foreignkey:CategoryId"`
	CategoryID uint     `gorm:"default:null"`

	Product   Category `gorm:"association_foreignkey:ProductId"`
	ProductID uint     `gorm:"default:null"`
}

//TagImages Scopes, not used
func TagImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "TagImage")
}

//CategoryImages ---
func CategoryImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "CategoryImage")
}

//ProductImages ---
func ProductImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "ProductImage")
}

// db.Scopes(CategoryImages, ProductImages).Find(&images)
