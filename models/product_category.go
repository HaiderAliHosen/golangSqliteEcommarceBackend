package models

//ProductCategory ---
type ProductCategory struct {
	Category   User `gorm:"association_foreignkey:CategoryId"`
	CategoryID uint
	Product    Product `gorm:"association_foreignkey:ProductId"`
	ProductID  uint
}

//TableName ---
func (*ProductCategory) TableName() string {
	return "products_categories"
}
