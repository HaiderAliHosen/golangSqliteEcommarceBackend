package models

//ProductTag ---
type ProductTag struct {
	Tag       User `gorm:"association_foreignkey:TagId"`
	TagID     uint
	Product   Product `gorm:"association_foreignkey:ProductId"`
	ProductID uint
}

//TableName ---
func (*ProductTag) TableName() string {
	return "products_tags"
}
