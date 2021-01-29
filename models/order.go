package models

import "github.com/jinzhu/gorm"

//Order ---
type Order struct {
	gorm.Model
	OrderStatus    int `gorm:"default:0"`
	TrackingNumber string

	OrderItems []OrderItem `gorm:"foreignKey:OrderId"`

	Address   Address `gorm:"association_foreignkey:AddressId:"`
	AddressID uint

	User            User `gorm:"foreignKey:UserId:"`
	UserID          uint `gorm:"default:null"`
	OrderItemsCount int  `gorm:"-"`
}

//GetOrderStatusAsString ---
func (order *Order) GetOrderStatusAsString() string {
	switch order.OrderStatus {
	case 0:
		return "processed"
	case 1:
		return "delivered"
	case 2:
		return "shipped"
	default:
		return "unknown"
	}
}
